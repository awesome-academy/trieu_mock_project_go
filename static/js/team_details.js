let currentOffset = 0;
const limit = 10;
let teamId = null;

$(document).ready(function () {
  if (!AuthService.isAuthenticated()) return;

  // Get team ID from URL
  const pathParts = window.location.pathname.split("/");
  teamId = pathParts[pathParts.length - 1];

  if (!teamId || isNaN(teamId)) {
    alert("Invalid Team ID");
    window.location.href = "/teams";
    return;
  }

  loadTeamDetails();
  loadTeamMembers(currentOffset);
});

/**
 * Fetch and display team details
 */
async function loadTeamDetails() {
  try {
    const team = await TeamService.getTeamDetails(teamId);
    $("#breadcrumb-team-name").text(team.name);
    $("#team-name-header").text(team.name);
    if (team.description) {
      $("#team-description").text(team.description);
    }

    // Update Team Info
    $("#team-leader").html(
      team.leader
        ? `<a href="/profile/${team.leader.id}" class="text-decoration-none">${team.leader.name}</a>`
        : "N/A"
    );
    $("#team-created-at").text(new Date(team.created_at).toLocaleString());
    $("#team-updated-at").text(new Date(team.updated_at).toLocaleString());

    // Update Projects
    updateProjectsList(team.projects);
  } catch (error) {
    console.error("Error fetching team details:", error);
    if (error.status !== 401) {
      alert("Failed to load team details.");
    }
  }
}

/**
 * Update projects list
 * @param {Array} projects
 */
function updateProjectsList(projects) {
  const container = $("#team-projects");
  if (!projects || projects.length === 0) {
    container.html('<p class="text-muted mb-0">No projects joined</p>');
    return;
  }

  let html = "";
  projects.forEach((project) => {
    html += `
      <div class="list-group-item px-0 border-0 border-bottom">
        <div class="d-flex w-100 justify-content-between">
          <h6 class="mb-1 fw-bold">${project.name}</h6>
          <small class="badge bg-secondary">${project.abbreviation}</small>
        </div>
        <small class="text-muted">
          ${
            project.start_date
              ? new Date(project.start_date).toLocaleDateString()
              : "N/A"
          } - 
          ${
            project.end_date
              ? new Date(project.end_date).toLocaleDateString()
              : "Present"
          }
        </small>
      </div>
    `;
  });
  container.html(html);
}

/**
 * Fetch and display team members
 * @param {number} offset
 */
async function loadTeamMembers(offset) {
  try {
    showLoadingTeamMembers();

    const response = await TeamService.getTeamMembers(teamId, limit, offset);
    updateMembersTable(response.members);
    updatePagination(response.page);
    currentOffset = offset;
  } catch (error) {
    console.error("Error fetching team members:", error);
    if (error.status !== 401) {
      alert("Failed to load team members.");
    }
  }
}

/**
 * Update members table with data
 * @param {Array} members
 */
function updateMembersTable(members) {
  const tbody = $("#members-table-body");
  if (!members || members.length === 0) {
    tbody.html(
      '<tr><td colspan="4" class="text-center text-muted">No members found in this team</td></tr>'
    );
    return;
  }

  let html = "";
  members.forEach((member) => {
    const joinedAt = member.joined_at
      ? new Date(member.joined_at).toLocaleDateString()
      : "N/A";
    html += `
      <tr class="member-row" onclick="goToUserProfile(${member.id})">
        <td>${member.id}</td>
        <td class="fw-bold text-primary">${member.name}</td>
        <td>${member.email}</td>
        <td>${joinedAt}</td>
      </tr>
    `;
  });
  tbody.html(html);
}

/**
 * Redirect to user profile
 * @param {number} userId
 */
function goToUserProfile(userId) {
  window.location.href = `/profile/${userId}`;
}

/**
 * Update pagination controls
 * @param {Object} pageInfo
 */
function updatePagination(pageInfo) {
  const pagination = $("#members-pagination");
  const totalPages = Math.ceil(pageInfo.total / pageInfo.limit);
  const currentPage = Math.floor(pageInfo.offset / pageInfo.limit) + 1;

  if (totalPages <= 1) {
    pagination.html("");
    return;
  }

  let html = "";

  // Previous button
  html += `
    <li class="page-item ${currentPage === 1 ? "disabled" : ""}">
      <a class="page-link" href="javascript:void(0)" onclick="loadTeamMembers(${
        (currentPage - 2) * limit
      })">Previous</a>
    </li>
  `;

  // Page numbers
  for (let i = 1; i <= totalPages; i++) {
    html += `
      <li class="page-item ${i === currentPage ? "active" : ""}">
        <a class="page-link" href="javascript:void(0)" onclick="loadTeamMembers(${
          (i - 1) * limit
        })">${i}</a>
      </li>
    `;
  }

  // Next button
  html += `
    <li class="page-item ${currentPage === totalPages ? "disabled" : ""}">
      <a class="page-link" href="javascript:void(0)" onclick="loadTeamMembers(${
        currentPage * limit
      })">Next</a>
    </li>
  `;

  pagination.html(html);
}

function showLoadingTeamMembers() {
  $("#members-table-body").html(`
      <tr>
        <td colspan="4" class="text-center">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </td>
      </tr>
    `);
}
