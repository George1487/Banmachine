using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Minio;

namespace Application;

public static class StorageExtensions
{

    public static IServiceCollection AddApplicationStorage(
        this IServiceCollection services
        , IConfiguration configuration)
    {
        services.Configure<StorageSettings>(configuration.GetSection("Minio"));

        var storageSettings = configuration.GetSection("Minio")
            .Get<StorageSettings>() ?? throw new InvalidOperationException("Minio section is not configured.");

        services.AddMinio(o => o
            .WithEndpoint(storageSettings.Endpoint)
            .WithCredentials(storageSettings.AccessKey, storageSettings.SecretKey)
            .WithSSL(false)
            .Build());
        return services;
    }
}
