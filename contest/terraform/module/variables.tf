## Required variables

variable "admins" {
  description = "Admins' GitHub IDs"
  type        = set(string)
}

variable "teams" {
  description = "List of teams and their members' GitHub IDs"
  type        = map(set(string))
  default = {
    "team1" = [
      "user1",
      "user2",
      "user3",
    ],
    "team2" = [
      "user4",
      "user5",
    ],
  }
}
# Example:
# teams = {
#   "team1" = [
#     "user1",
#     "user2",
#     "user3",
#   ],
#   "team2" = [
#     "user4",
#     "user5",
#   ],
# }

## Optional variables

variable "name" {
  description = "Name for the main resources. Useful for creating test resources avoiding name conflicts"
  type        = string
  default     = "ishocon1"
}

variable "ami_id" {
  description = "AMI ID for the ISHOCON EC2 instances"
  type        = string
  default     = "ami-05861adff828ed920"

}

variable "instance_type" {
  description = "Instance type for the ISHOCON EC2 instances"
  type        = string
  default     = "c7i.xlarge"

}

variable "use_spot_instance" {
  description = "Use spot EC2 instance to save cost"
  type        = bool
  default     = false
}

variable "vpc_cidr_block" {
  description = "VPC CIDR block for the ISHOCON VPC"
  type        = string
  default     = "172.16.0.0/16"
}


