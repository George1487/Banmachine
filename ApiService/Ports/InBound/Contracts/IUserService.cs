using Domain.Users;
using Ports.InBound.DTO;

namespace Ports.InBound.Contracts;

public interface IUserService
{
    
    UserResult RegisterUser(User user);
    
    UserResult GetUser(Guid userId);
    
    LoginResult Login(string email, string password);
}