FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

COPY ["ApiService.sln", "./"]
COPY ["Directory.build.props", "Directory.Build.props"]
COPY ["Application/Application.csproj", "Application/"]
COPY ["Domain/Domain.csproj", "Domain/"]
COPY ["Infrastructure/Infrastructure.csproj", "Infrastructure/"]
COPY ["Ports/Ports.csproj", "Ports/"]
COPY ["Presentation/Presentation.csproj", "Presentation/"]
COPY ["Web/Web.csproj", "Web/"]

RUN dotnet restore "Web/Web.csproj"

COPY . .
RUN dotnet publish "Web/Web.csproj" -c Release -o /app/publish /p:UseAppHost=false

FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS final
WORKDIR /app

ENV ASPNETCORE_URLS=http://+:8080
EXPOSE 8080

COPY --from=build /app/publish .

ENTRYPOINT ["dotnet", "Web.dll"]
