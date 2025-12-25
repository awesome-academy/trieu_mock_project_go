$(document).ready(function () {
  const token = getAccessToken();
  if (!token) return;

  $.ajax({
    url: "/api/profile",
    method: "GET",
    headers: {
      Authorization: "Bearer " + token,
    },
    success: function (data) {
      // Update Header
      $("#profile-name-header").text(data.name);
      if (data.position) {
        $("#profile-position-header").text(data.position.name);
      } else {
        $("#profile-position-header").text("No Position");
      }
      if (data.current_team) {
        $("#profile-team-header").text(data.current_team.name);
      } else {
        $("#profile-team-header").text("No Team");
      }

      // Update Basic Info
      $("#info-name").text(data.name);
      $("#info-email").text(data.email);

      if (data.birthday) {
        const birthday = new Date(data.birthday);
        $("#info-birthday").text(birthday.toLocaleDateString());
      } else {
        $("#info-birthday").text("N/A");
      }

      if (data.current_team) {
        $("#info-team").text(data.current_team.name);
      } else {
        $("#info-team").text("N/A");
      }

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
    },
    error: function (xhr) {
      console.error("Error fetching profile:", xhr);
      if (xhr.status === 401) {
        // auth.js handles redirect if token is invalid,
        // but we can also trigger logout here if the API returns 401
        logout();
      } else {
        alert("Failed to load profile information. Please try again later.");
      }
    },
  });
});
