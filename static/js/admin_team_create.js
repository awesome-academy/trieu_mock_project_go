document.addEventListener("DOMContentLoaded", function () {
  const leaderSearch = document.getElementById("leaderSearch");
  const leaderSearchResults = document.getElementById("leaderSearchResults");
  const leaderID = document.getElementById("leaderID");
  const selectedLeader = document.getElementById("selectedLeader");
  const selectedLeaderName = document.getElementById("selectedLeaderName");
  const clearLeader = document.getElementById("clearLeader");
  const createTeamBtn = document.getElementById("createTeamBtn");
  const createTeamForm = document.getElementById("createTeamForm");

  let searchTimeout;

  leaderSearch.addEventListener("focus", function () {
    searchUsers(this.value.trim());
  });

  leaderSearch.addEventListener("input", function () {
    clearTimeout(searchTimeout);
    const query = this.value.trim();
    searchTimeout = setTimeout(() => searchUsers(query), 300);
  });

  async function searchUsers(query) {
    try {
      const response = await fetch(
        `/admin/users/search?name=${encodeURIComponent(query)}&limit=5&offset=0`
      );
      const data = await response.json();

      leaderSearchResults.innerHTML = "";
      if (data.users && data.users.length > 0) {
        data.users.forEach((user) => {
          const item = document.createElement("button");
          item.type = "button";
          item.className = "list-group-item list-group-item-action";
          const teamInfo = user.current_team
            ? ` [Team: ${user.current_team.name}]`
            : " [No Team]";
          item.textContent = `${user.name} (${user.email})${teamInfo}`;
          item.addEventListener("click", () => selectLeader(user));
          leaderSearchResults.appendChild(item);
        });
        leaderSearchResults.style.display = "block";
      } else {
        leaderSearchResults.innerHTML =
          '<div class="list-group-item">No users found</div>';
        leaderSearchResults.style.display = "block";
      }
    } catch (error) {
      console.error("Search error:", error);
    }
  }

  function selectLeader(user) {
    if (user.current_team) {
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

  // Close search results when clicking outside
  document.addEventListener("click", (e) => {
    if (
      !leaderSearch.contains(e.target) &&
      !leaderSearchResults.contains(e.target)
    ) {
      leaderSearchResults.style.display = "none";
    }
  });

  createTeamBtn.addEventListener("click", async () => {
    if (!createTeamForm.checkValidity()) {
      createTeamForm.reportValidity();
      return;
    }

    const formData = new FormData(createTeamForm);
    const data = {
      name: formData.get("name"),
      description: formData.get("description"),
      leader_id: parseInt(formData.get("leader_id")),
    };

    try {
      await AdminTeamService.createTeam(data);
      window.location.href = "/admin/teams";
    } catch (error) {
      alert(error.message);
    }
  });
});
