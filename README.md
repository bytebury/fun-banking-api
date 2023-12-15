# Golfer ⛳️

This is a skeleton project for how you _could_ setup a go/Gin server with JWT and Database boilerplate. Docker configurations are also shipped in the box.

## Getting Started

I tried to make it as easy as possible to create a brand new Go API with basic features. This is what I came up with.

## Sample `.env` File

```env
APP_BASE_URL=http://localhost:8080
DATABASE_URL=postgresql://username:password@localhost:5432/database
JWT_SECRET=TOP_SECRET_K3Y

EMAIL_USERNAME=your.email@domain.com
EMAIL_PASSWORD=top_secret_sauce
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
```

You will also need to update application information in the `config/constants.go` file to customize this project to your needs.

## Building

### Local

```
go run main.go
```

### Docker

```
docker build -t golfer-dev .
docker run -p 8080:8080 golfer-dev
```
