// FIXME
# locals {
#   admins = toset([
#     "FIXME_github_id_for_admin",
#   ])

#   players = toset([
#     "FIXME_github_id_for_player1",
#     "FIXME_github_id_for_player2",
#   ])
#   # TODO: support team format
# }

locals {
  admins = toset([
    "showwin",
  ])

  players = toset([
    "showwin",
  ])
}

module "main" {
  source = "../module"

  admins = local.admins
  players = local.players

  use_spot_instance = true
}
