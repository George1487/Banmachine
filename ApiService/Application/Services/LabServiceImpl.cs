using Domain.Labs;
using Domain.Users;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;

namespace Application.Services;

public class LabServiceImpl : ILabService
{
    
    private readonly ILabRepository  _labRepo;

    public LabServiceImpl(ILabRepository labRepo)
    {
        _labRepo = labRepo;
    }
    
    public LabResult GetLab(Guid labId)
    {
        return _labRepo.GetLab(labId);
    }

    public LabsResult GetLabs(Guid userId, UserRole role)
    {
        var result = _labRepo.GetLabs();
        if (result is LabsResult.Failure failure)
        {
            return new LabsResult.Failure(failure.Reason);
        }
        
        var success = (LabsResult.Success)result;
        if (role == UserRole.Teacher)
        {
            var lab = new LabsResult.Success(
                success.Labs.Where(o => o.TeacherId == userId).ToList());
            return lab;
        }

        return success;
    }

    public LabResult AddLab(Lab lab)
    {
        if (lab.DeadlineAt.ToUniversalTime() < DateTime.UtcNow)
        {
            return new LabResult.Failure("Deadline must be in the future");
        }
        return _labRepo.AddLab(lab);
    }

    public LabResult PatchLab(
        Guid labId,
        Guid userId, 
        LabStatus newLabStatus, 
        string newTitle, 
        string newDescription,
        DateTimeOffset newDeadlineAt)
    {
        
        var lab = _labRepo.GetLab(labId);
        if (lab is LabResult.Failure failure)
        {
            return new LabResult.Failure(failure.Reason);
        }
        var success = (LabResult.Success)lab;
        if (userId != success.Lab.TeacherId)
        {
            return new LabResult.Failure("Not allowed");
        }

        return _labRepo.PatchLab(
            labId,
            newLabStatus,
            newTitle,
            newDescription,
            newDeadlineAt);
    }
}
