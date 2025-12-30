let currentOffset = 0;
const limit = 10;

$(document).ready(function () {
  if (!AuthService.isAuthenticated()) return;

  loadTeams(currentOffset);
});

/**
 * Fetch and display teams
 * @param {number} offset
 */
async function loadTeams(offset) {
  try {
    showLoadingTeams();

    const response = await TeamService.listTeams(limit, offset);
    updateTeamsTable(response.teams);
    updatePagination(response.page);
    currentOffset = offset;
  } catch (error) {
    console.error("Error fetching teams:", error);
    if (error.status !== 401) {
      alert("Failed to load teams. Please try again later.");
    }
  }
}

/**
 * Update teams table with data
 * @param {Array} teams
 */
function updateTeamsTable(teams) {
  const tbody = $("#teams-table-body");
  if (!teams || teams.length === 0) {
    tbody.html(
      '<tr><td colspan="6" class="text-center text-muted">No teams found</td></tr>'
    );
    return;
  }

  let html = "";
  teams.forEach((team) => {
    const createdAt = new Date(team.created_at).toLocaleDateString();
    html += `
      <tr class="team-row" onclick="goToTeamDetails(${team.id})">
        <td>${team.id}</td>
        <td class="fw-bold">${team.name}</td>
        <td>${team.leader ? team.leader.name : "N/A"}</td>
        <td><span class="badge bg-info text-dark">${
          team.members ? team.members.length : 0
        } Members</span></td>
        <td><span class="badge bg-secondary">${
          team.projects ? team.projects.length : 0
        } Projects</span></td>
        <td>${createdAt}</td>
      </tr>
    `;
  });
  tbody.html(html);
}

/**
 * Redirect to team details page
 * @param {number} teamId
 */
function goToTeamDetails(teamId) {
  window.location.href = `/teams/${teamId}`;
}

/**
 * Update pagination controls
 * @param {Object} pageInfo
 */
function updatePagination(pageInfo) {
  const pagination = $("#teams-pagination");
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
      <a class="page-item page-link" href="javascript:void(0)" onclick="loadTeams(${
        (currentPage - 2) * limit
      })">Previous</a>
    </li>
  `;

  // Page numbers
  for (let i = 1; i <= totalPages; i++) {
    html += `
      <li class="page-item ${i === currentPage ? "active" : ""}">
        <a class="page-link" href="javascript:void(0)" onclick="loadTeams(${
          (i - 1) * limit
        })">${i}</a>
      </li>
    `;
  }

  // Next button
  html += `
    <li class="page-item ${currentPage === totalPages ? "disabled" : ""}">
      <a class="page-link" href="javascript:void(0)" onclick="loadTeams(${
        currentPage * limit
      })">Next</a>
    </li>
  `;

  pagination.html(html);
}
function showLoadingTeams() {
  $("#teams-table-body").html(`
      <tr>
        <td colspan="6" class="text-center">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </td>
      </tr>
    `);
}
