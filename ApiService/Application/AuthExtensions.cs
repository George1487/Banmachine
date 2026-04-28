using System.Text;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.IdentityModel.Tokens;

namespace Application;

public static class AuthExtensions
{

    public static IServiceCollection AddJwtAuthentication(
        this IServiceCollection services, IConfiguration configuration)
    {
        services.Configure<AuthSettings>(configuration.GetSection("AuthSettings"));

        var authSettings = configuration.GetSection("AuthSettings").Get<AuthSettings>()
            ?? throw new InvalidOperationException("AuthSettings section is not configured.");

        if (string.IsNullOrWhiteSpace(authSettings.SecretKey))
        {
            throw new InvalidOperationException("AuthSettings:SecretKey is not configured.");
        }

        services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
            .AddJwtBearer(o =>
            {
                o.TokenValidationParameters = new TokenValidationParameters
                {
                    ValidateIssuer = false,
                    ValidateAudience = false,
                    ValidateLifetime = true,
                    ValidateIssuerSigningKey = true,
                    IssuerSigningKey = new SymmetricSecurityKey(
                        Encoding.UTF8.GetBytes(authSettings!.SecretKey)),
                };
            });
        
        return services;
    }
    
}
