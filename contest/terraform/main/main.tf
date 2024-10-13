// FIXME
locals {
  admins = [
    "admin_github_id",
  ]

  teams = {
    "team1" = [
      "user1_github_id",
      "user2",
      "user3",
    ],
    "team2" = [
      "user4",
      "user5",
    ],
  }
}

module "main" {
  source = "../module"

  admins = local.admins
  teams = local.teams

  use_spot_instance = false
}
