using Microsoft.Extensions.DependencyInjection;

namespace Presentation;

public static class ServiceCollectionExtensions
{

    public static IServiceCollection AddPresentation(this IServiceCollection services)
    {
        services.AddControllers();
        return services;
    }
    
    
}