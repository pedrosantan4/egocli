// cmd/templates.go
package cmd

// ModuleTemplate define a estrutura dos templates
type ModuleTemplate struct {
	DirName     string
	FileName    string
	Content     string
	CommandType string
}

// templates é o mapa global de templates AWS
var Templates = map[string]ModuleTemplate{
	"vpc": {
		DirName:     "01-networking",
		FileName:    "vpc.tf",
		CommandType: "gen",
		Content: `resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = {
    Name        = "egocli-vpc"
    Environment = "dev"
  }
}

resource "aws_subnet" "public" {
  count                   = 2
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.${count.index + 1}.0/24"
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  
  tags = {
    Name = "public-subnet-${count.index + 1}"
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}`,
	},

	"eks": {
		DirName:     "02-kubernetes",
		FileName:    "eks.tf",
		CommandType: "gen",
		Content: `resource "aws_eks_cluster" "main" {
  name     = "egocli-cluster"
  role_arn = aws_iam_role.eks_cluster.arn
  version  = "1.27"

  vpc_config {
    subnet_ids = aws_subnet.public[*].id
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_cluster_policy,
  ]
}

resource "aws_iam_role" "eks_cluster" {
  name = "eks-cluster-role"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster.name
}`,
	},

	"rds": {
		DirName:     "03-database",
		FileName:    "rds.tf",
		CommandType: "gen",
		Content: `resource "aws_db_instance" "main" {
  identifier = "egocli-db"
  
  engine         = "postgres"
  engine_version = "14.9"
  instance_class = "db.t3.micro"
  
  allocated_storage = 20
  storage_type      = "gp2"
  
  db_name  = "egocli"
  username = "postgres"
  password = "changeme123"
  
  vpc_security_group_ids = [aws_security_group.rds.id]
  
  skip_final_snapshot = true
  
  tags = {
    Name = "egocli-database"
  }
}

resource "aws_security_group" "rds" {
  name_prefix = "rds-"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }
}`,
	},

	"s3": {
		DirName:     "04-storage",
		FileName:    "s3.tf",
		CommandType: "gen",
		Content: `resource "aws_s3_bucket" "main" {
  bucket = "egocli-bucket-${random_string.suffix.result}"
  
  tags = {
    Name = "egocli-storage"
  }
}

resource "aws_s3_bucket_versioning" "main" {
  bucket = aws_s3_bucket.main.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "main" {
  bucket = aws_s3_bucket.main.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}`,
	},

	"iam": {
		DirName:     "05-security",
		FileName:    "iam.tf",
		CommandType: "gen",
		Content: `resource "aws_iam_role" "app_role" {
  name = "egocli-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_policy" "app_policy" {
  name = "egocli-app-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:PutObject",
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
      Resource = "*"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "app_policy_attachment" {
  role       = aws_iam_role.app_role.name
  policy_arn = aws_iam_policy.app_policy.arn
}`,
	},

	"lambda": {
		DirName:     "06-functions",
		FileName:    "lambda.tf",
		CommandType: "gen",
		Content: `resource "aws_lambda_function" "main" {
  filename         = "lambda.zip"
  function_name    = "egocli-function"
  role            = aws_iam_role.lambda_role.arn
  handler         = "index.handler"
  runtime         = "nodejs18.x"
  
  tags = {
    Name = "egocli-lambda"
  }
}

resource "aws_iam_role" "lambda_role" {
  name = "lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}`,
	},
}
