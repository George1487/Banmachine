using Microsoft.AspNetCore.Http;

namespace Presentation.RequestDTO;

public sealed class CreateSubmissionRequest
{
    public IFormFile File { get; set; } = null!;
}
