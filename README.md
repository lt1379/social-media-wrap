# Go Project

## Create migration
```aiignore
liquibase init project --project-dir=my_project_postgres --changelog-file=my-project-changelog.sql --format=sql --project-defaults-file=liquibase.properties --url=jdbc:postgresql://localhost:5432/my_project --username=project --password=MyPassword_123
```

## Update migration
```aiignore
liquibase update --changelog-file=my-project-changelog.sql --url=jdbc:postgresql://localhost:5432/my_project --username=project --password=MyPassword_123
```