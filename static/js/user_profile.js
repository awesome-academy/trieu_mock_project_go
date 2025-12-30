$(document).ready(function () {
  if (!AuthService.isAuthenticated()) return;

  loadUserProfile();
});

/**
 * Fetch and display user profile
 */
async function loadUserProfile() {
  try {
    const data = await UserService.getProfile();
    updateProfileDOM(data);
  } catch (error) {
    console.error("Error fetching profile:", error);
    // API utility handles 401, so we only handle other errors here
    if (error.status !== 401) {
      alert("Failed to load profile information. Please try again later.");
    }
  }
}

/**
 * Update DOM with profile data
 * @param {Object} data
 */
function updateProfileDOM(data) {
  // Update Header
  $("#profile-name-header").text(data.name);
  $("#profile-position-header").text(
    data.position ? data.position.name : "No Position"
  );
  $("#profile-team-header").text(
    data.current_team ? data.current_team.name : "No Team"
  );

  // Update Basic Info
  $("#info-name").text(data.name);
  $("#info-email").text(data.email);

  if (data.birthday) {
    const birthday = new Date(data.birthday);
    $("#info-birthday").text(birthday.toLocaleDateString());
  } else {
    $("#info-birthday").text("N/A");
  }

  $("#info-team").text(data.current_team ? data.current_team.name : "N/A");

  if (data.position) {
    $("#info-position").text(
      `${data.position.name} (${data.position.abbreviation})`
    );
  } else {
    $("#info-position").text("N/A");
  }

  // Update Skills
  if (data.skills && data.skills.length > 0) {
    let skillsHtml = "";
    data.skills.forEach((skill) => {
      skillsHtml += `<span class="badge bg-primary badge-skill">${skill.name}</span>`;
    });
    $("#profile-skills").html(skillsHtml);
  } else {
    $("#profile-skills").html(
      '<p class="text-muted mb-0">No skills listed</p>'
    );
  }

  // Update Projects
  if (data.projects && data.projects.length > 0) {
    let projectsHtml = "";
    data.projects.forEach((project) => {
      const startDate = project.start_date
        ? new Date(project.start_date).toLocaleDateString()
        : "N/A";
      const endDate = project.end_date
        ? new Date(project.end_date).toLocaleDateString()
        : "Present";
      projectsHtml += `
        <tr>
          <td class="fw-bold">${project.name}</td>
          <td><span class="badge bg-secondary">${project.abbreviation}</span></td>
          <td>${startDate}</td>
          <td>${endDate}</td>
        </tr>
      `;
    });
    $("#profile-projects").html(projectsHtml);
  } else {
    $("#profile-projects").html(
      '<tr><td colspan="4" class="text-center text-muted">No projects joined</td></tr>'
    );
  }

  // Update Avatar with name
  const avatarUrl = `https://ui-avatars.com/api/?name=${encodeURIComponent(
    data.name
  )}&background=random&size=150`;
  $('img[alt="avatar"]').attr("src", avatarUrl);
}
