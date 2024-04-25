# Technology selection

## Backend

### Language: Golang
The Golang is statically typed and has strict compile-time error checking, making it suitable for large team development. It is easy to clearly separate responsibilities for each package and promotes code standardization. Furthermore, Go does not have exception handling, and is designed to require explicit error handling each time an error occurs, which can reduce the occurrence of unexpected bugs. It also enables Go's parallel processing and fast execution speed. For these reasons, I chose my favorite language, Go.

### Framework: Echo
Among the Golang web frameworks that are popular in Japan, Gin and Echo are particularly popular. I chose Echo because I already had work experience with Gin and was willing to try new technology. Echo's selling points are high performance, scalability, and minimalism. In fact, I expanded the middleware to store Access Token in Redis.

### O/R Mapper: GORM
There are many O/R mappers for Go lang, but among the major O/R mappers, we decided to use GORM, which had the most stars on Github at that time.

## Frontend 

### Language: Typescript
Because Typescript >>> Javascript in front-end development. It was adopted with considerations such as type safety, code completion, and readability.

### Framework: Next.js(SSG)
As the service expanded, a more complex UI and component reuse became necessary. By providing a component-based architecture, Next.js can greatly improve UI reusability and maintainability. In addition, Next.js' static site generation (SSG) feature generates HTML, CSS, and JS at build time and delivers them to Echo servers. This reduces the strain on server resources while providing end users with faster page loading speeds.

### CSS framework: PandaCSS
PandaCSS is a CSS-in-JS library that, in combination with TypeScript, ensures type safety for CSS. This combination allows developers to benefit from code completion and compile-time type checking, making it easier to proactively detect style-related errors. Additionally, PandaCSS only generates minimal CSS at build time, making the final file size very lightweight. I chose PandaCSS because of its performance optimization and high developer experience.

## Dependent service
The reason for selecting external services is that they can basically be used within the free plan.

### Deployment platform: Render
Render's free plan offers up to 750 hours of service operation per month and supports deployment using Docker images. Additionally, it facilitates automatic deployment when changes are merged into the main branch.

[Render pricing](https://render.com/pricing)

### Database and storage: Supabase
Database: Full access to PostgreSQL databases, and the free plan provides 500MB of database storage per month.
Storage: Includes storage features that support file uploads and downloads, with the initial plan offering 2GB of storage and 1GB of network transfer per month.

[Supabase pricing](https://supabase.com/pricing)

### Session store (Redis): Upstash
Upstash's free plan provides 30MB of data storage, supports up to 20 connections at the same time, and can process up to 10 requests per second, meeting the connectivity requirements needed for small-scale applications like personal development.

[Upstash pricing](https://upstash.com/pricing)

### Monitoring: UptimeRobot
I leverage UptimeRobot to monitor my application's health check endpoint every 5 minutes. This effectively prevents the 15 minute idle timeout due to Render's auto-sleep feature and the 1 week idle timeout of Supabase's PostgreSQL. As a result, I continually check the health of our applications and ensure that my service is always up and running.

[UptimeRobot pricing](https://uptimerobot.com/pricing/)
