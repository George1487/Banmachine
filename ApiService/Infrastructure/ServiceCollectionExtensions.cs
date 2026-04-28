using Infrastructure.Adapters;
using Infrastructure.Repositories;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Ports.OutBound.Contracts;

namespace Infrastructure;

public static class ServiceCollectionExtensions
{
    public static IServiceCollection AddInfrastructure(
        this IServiceCollection services,
        IConfiguration configuration)
    {
        var connectionString = configuration.GetConnectionString("DefaultConnection");
        if (string.IsNullOrWhiteSpace(connectionString))
        {
            throw new InvalidOperationException(
                "ConnectionStrings:DefaultConnection is not configured.");
        }

        services.AddDbContext<AppDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<AnalysisRepository>();
        services.AddScoped<UserRepository>();
        services.AddScoped<LabRepository>();
        services.AddScoped<SubmissionRepository>();
        services.AddScoped<IngestJobRepository>();
        services.AddScoped<SubmissionAnalysisSummaryRepository>();
        services.AddScoped<PairwiseSimilarityRepository>();

        services.AddScoped<IAnalysisRepository, AnalysisRepositoryAdapter>();
        services.AddScoped<IUserRepository, UserRepositoryAdapter>();
        services.AddScoped<ILabRepository, LabRepositoryAdapter>();
        services.AddScoped<ISubmissionRepository, SubmissionRepositoryAdapter>();
        services.AddScoped<IIngestJobRepository, IngestJobRepositoryAdapter>();
        services.AddScoped<ISubmissionAnalysisSummaryRepository, SubmissionAnalysisSummaryRepositoryAdapter>();
        services.AddScoped<IPairwiseSimilarityRepository, PairwiseSimilarityRepositoryAdapter>();

        return services;
    }
}
