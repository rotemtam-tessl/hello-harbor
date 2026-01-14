data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./models",
    "--dialect", "sqlite", // | postgres | sqlite | sqlserver
  ]
}
env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://dev?mode=memory"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}