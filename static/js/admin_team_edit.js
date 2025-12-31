document.addEventListener("DOMContentLoaded", function () {
  const teamId = document.getElementById("teamId").value;
  const leaderSearch = document.getElementById("leaderSearch");
  const leaderSearchResults = document.getElementById("leaderSearchResults");
  const leaderID = document.getElementById("leaderID");
  const selectedLeader = document.getElementById("selectedLeader");
  const selectedLeaderName = document.getElementById("selectedLeaderName");
  const clearLeader = document.getElementById("clearLeader");
  const updateTeamBtn = document.getElementById("updateTeamBtn");
  const editTeamForm = document.getElementById("editTeamForm");

  const memberSearch = document.getElementById("memberSearch");
  const memberSearchResults = document.getElementById("memberSearchResults");
  const memberList = document.getElementById("memberList");
  const noMembersRow = document.getElementById("noMembersRow");
  const historyContainer = document.getElementById("historyContainer");

  let searchTimeout;

  // History Pagination and Loading
  async function loadHistory(offset = 0) {
    try {
      const response = await fetch(
        `/admin/teams/${teamId}/history/partial?offset=${offset}&limit=10`
      );
      const html = await response.text();
      historyContainer.innerHTML = html;
      attachHistoryPaginationEvents();
    } catch (error) {
      console.error("Failed to load history:", error);
    }
  }

  function attachHistoryPaginationEvents() {
    const paginationLinks =
      historyContainer.querySelectorAll(".history-page-link");
    paginationLinks.forEach((link) => {
      link.addEventListener("click", (e) => {
        e.preventDefault();
        const offset = e.target.getAttribute("data-offset");
        loadHistory(offset);
      });
    });
  }

  attachHistoryPaginationEvents();

  // Leader Search Logic
  leaderSearch.addEventListener("focus", function () {
    searchUsers(this.value.trim(), leaderSearchResults, selectLeader);
  });

  leaderSearch.addEventListener("input", function () {
    clearTimeout(searchTimeout);
    const query = this.value.trim();
    searchTimeout = setTimeout(
      () => searchUsers(query, leaderSearchResults, selectLeader),
      300
    );
  });

  function selectLeader(user) {
    if (user.current_team && user.current_team.id != teamId) {
      alert(
        "Leader already belongs to a team. Please remove them from their current team or choose another leader."
      );
      return;
    }
    leaderID.value = user.id;
    selectedLeaderName.textContent = user.name;
    selectedLeader.classList.remove("d-none");
    leaderSearch.value = "";
    leaderSearchResults.style.display = "none";
    leaderSearch.disabled = true;
  }

  clearLeader.addEventListener("click", () => {
    leaderID.value = "";
    selectedLeader.classList.add("d-none");
    leaderSearch.disabled = false;
    leaderSearch.focus();
  });

  // Member Search Logic
  memberSearch.addEventListener("focus", function () {
    searchUsers(this.value.trim(), memberSearchResults, addMember);
  });

  memberSearch.addEventListener("input", function () {
    clearTimeout(searchTimeout);
    const query = this.value.trim();
    searchTimeout = setTimeout(
      () => searchUsers(query, memberSearchResults, addMember),
      300
    );
  });

  async function searchUsers(query, resultsContainer, onSelect) {
    try {
      const response = await fetch(
        `/admin/users/search?name=${encodeURIComponent(query)}&limit=5&offset=0`
      );
      const data = await response.json();

      resultsContainer.innerHTML = "";
      if (data.users && data.users.length > 0) {
        data.users.forEach((user) => {
          const item = document.createElement("button");
          item.type = "button";
          item.className = "list-group-item list-group-item-action";
          const teamInfo = user.current_team
            ? ` [Team: ${user.current_team.name}]`
            : " [No Team]";
          item.textContent = `${user.name} (${user.email})${teamInfo}`;
          item.addEventListener("click", () => onSelect(user));
          resultsContainer.appendChild(item);
        });
        resultsContainer.style.display = "block";
      } else {
        resultsContainer.innerHTML =
          '<div class="list-group-item">No users found</div>';
        resultsContainer.style.display = "block";
      }
    } catch (error) {
      console.error("Search error:", error);
    }
  }

  async function addMember(user) {
    if (user.current_team) {
      if (user.current_team.id == teamId) {
        memberSearch.value = "";
        memberSearchResults.style.display = "none";
        return;
      }
      const currentTeamName = document.getElementById("name").value;
      if (
        !confirm(
          `User is currently in team "${user.current_team.name}". Do you want to move them to "${currentTeamName}"?`
        )
      ) {
        memberSearch.value = "";
        memberSearchResults.style.display = "none";
        return;
      }
    }
    try {
      await AdminTeamService.addMember(teamId, user.id);

      // Reload history immediately to reflect the change
      loadHistory();

      // Update UI
      if (noMembersRow) noMembersRow.remove();

      // Check if user already in list
      if (document.querySelector(`tr[data-user-id="${user.id}"]`)) {
        memberSearch.value = "";
        memberSearchResults.style.display = "none";
        return;
      }

      const row = document.createElement("tr");
      row.setAttribute("data-user-id", user.id);
      row.innerHTML = `
        <td>${user.name}</td>
        <td>
          <button class="btn btn-sm btn-outline-danger remove-member-btn" data-user-id="${user.id}">
            <i class="bi bi-trash"></i>
          </button>
        </td>
      `;
      memberList.appendChild(row);
      attachRemoveEvent(row.querySelector(".remove-member-btn"));

      memberSearch.value = "";
      memberSearchResults.style.display = "none";
    } catch (error) {
      alert(error.message);
    }
  }

  function attachRemoveEvent(btn) {
    btn.addEventListener("click", async function () {
      const userId = this.getAttribute("data-user-id");
      if (confirm("Are you sure you want to remove this member?")) {
        try {
          await AdminTeamService.removeMember(teamId, userId);
          this.closest("tr").remove();
          if (memberList.children.length === 0) {
            memberList.innerHTML =
              '<tr id="noMembersRow"><td colspan="2" class="text-center">No members in this team</td></tr>';
          }
          // Reload history
          loadHistory();
        } catch (error) {
          alert(error.message);
        }
      }
    });
  }

  document.querySelectorAll(".remove-member-btn").forEach(attachRemoveEvent);

  // Close search results when clicking outside
  document.addEventListener("click", (e) => {
    if (
      !leaderSearch.contains(e.target) &&
      !leaderSearchResults.contains(e.target)
    ) {
      leaderSearchResults.style.display = "none";
    }
    if (
      !memberSearch.contains(e.target) &&
      !memberSearchResults.contains(e.target)
    ) {
      memberSearchResults.style.display = "none";
    }
  });

  updateTeamBtn.addEventListener("click", async () => {
    if (!editTeamForm.checkValidity()) {
      editTeamForm.reportValidity();
      return;
    }

    const formData = new FormData(editTeamForm);
    const data = {
      name: formData.get("name"),
      description: formData.get("description"),
      leader_id: parseInt(formData.get("leader_id")),
    };

    try {
      await AdminTeamService.updateTeam(teamId, data);
      window.location.href = "/admin/teams";
    } catch (error) {
      alert(error.message);
    }
  });

  // Initial state for leader search if already selected
  if (leaderID.value) {
    leaderSearch.disabled = true;
  }
});
