using Application.Services;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Ports.InBound.Contracts;

namespace Application;

public static class ServiceCollectionExtensions
{
    public static IServiceCollection AddApplicationServices(
        this IServiceCollection services,
        IConfiguration configuration)
    {
        services.AddScoped<ITokenService, TokenServiceImpl>();
        services.AddScoped<IUserService, UserServiceImpl>();
        services.AddScoped<ILabService, LabServiceImpl>();
        services.AddScoped<IAnalysisService, AnalysisJobServiceImpl>();
        services.AddScoped<IPairwiseSimilarityService, PairwiseSimilarityServiceImpl>();
        services.AddScoped<ISubmissionAnalysisSummaryService, 
            SubmissionAnalysisSummaryServiceImpl>();
        services.AddScoped<ISubmissionService, SubmissionServiceImpl>();
        return services;
    }
}
