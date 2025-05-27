env "local" {
  url = "mysql://root:password@localhost:3306/im_db"
  dev = "mysql://root:password@localhost:3306/im_db"
  migration {
    dir = "file://migrations"
  }
} 