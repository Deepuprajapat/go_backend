env "local" {
  url = "mysql://root:password@localhost:3306/im_db_dev"
  dev = "mysql://root:password@localhost:3306/im_db_dev"
  migration {
    dir = "file://migrations"
  }
} 