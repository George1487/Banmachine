using Domain.Users;

namespace Ports.InBound.Contracts;

public interface ITokenService
{
    
    string CreateToken(User user);
    
    Guid GetIdFromToken(string token);
}