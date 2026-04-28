using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using Domain.Users;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.Tokens;
using Ports.InBound.Contracts;

namespace Application.Services;

public class TokenServiceImpl(IOptions<AuthSettings> options) : ITokenService
{
    public string CreateToken(User user)
    {

        var claims = new List<Claim>
        {
            new Claim("UserId", user.UserId.ToString()),
            new Claim(ClaimTypes.Role, user.Role.ToString())
        };

        var token = new JwtSecurityToken(
            expires: DateTime.UtcNow.AddHours(2),
            claims: claims,
            signingCredentials: new SigningCredentials(
                new SymmetricSecurityKey(
               Encoding.UTF8.GetBytes(options.Value.SecretKey)),
               SecurityAlgorithms.HmacSha256)
            );
        
        return new JwtSecurityTokenHandler().WriteToken(token);
    }

    public Guid GetIdFromToken(string token)
    {
        throw new NotImplementedException();
    }
}