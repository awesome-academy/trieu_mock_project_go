document.addEventListener("DOMContentLoaded", function () {
  const createProjectForm = document.getElementById("createProjectForm");
  const teamSelect = document.getElementById("team");
  const memberList = document.getElementById("memberList");
  const memberPlaceholder = document.getElementById("memberPlaceholder");
  const memberCheckboxList = document.getElementById("memberCheckboxList");
  const selectedMembersContainer = document.getElementById(
    "selectedMembersContainer"
  );
  const noMembersSelected = document.getElementById("noMembersSelected");
  const createProjectBtn = document.getElementById("createProjectBtn");

  const leaderSearch = document.getElementById("leaderSearch");
  const leaderSearchResults = document.getElementById("leaderSearchResults");
  const leaderID = document.getElementById("leaderID");
  const selectedLeader = document.getElementById("selectedLeader");
  const selectedLeaderName = document.getElementById("selectedLeaderName");
  const clearLeader = document.getElementById("clearLeader");

  let currentTeamId = "";
  let searchTimeout;
  let selectedMembers = new Map(); // Map<id, name>

  // Handle leader search
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
          item.textContent = `${user.name} (${user.email})`;
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
      leaderSearch &&
      !leaderSearch.contains(e.target) &&
      leaderSearchResults &&
      !leaderSearchResults.contains(e.target)
    ) {
      leaderSearchResults.style.display = "none";
    }
  });

  // Handle team change
  teamSelect.addEventListener("change", async function () {
    const newTeamId = this.value;

    if (selectedMembers.size > 0 && currentTeamId !== "") {
      if (
        !confirm(
          "Changing the team will clear all currently selected members. Do you want to continue?"
        )
      ) {
        // Revert team selection
        teamSelect.value = currentTeamId;
        return;
      }
    }

    currentTeamId = newTeamId;
    selectedMembers.clear();
    updateSelectedMembersUI();
    await loadTeamMembers(newTeamId);
  });

  async function loadTeamMembers(teamId) {
    if (!teamId) {
      memberList.classList.add("d-none");
      memberPlaceholder.classList.remove("d-none");
      return;
    }

    try {
      const members = await AdminProjectService.getTeamMembers(teamId);

      memberCheckboxList.innerHTML = "";
      if (members && members.length > 0) {
        members.forEach((member) => {
          const div = document.createElement("div");
          div.className = "px-3 py-1 hover-bg-light";

          const formCheck = document.createElement("div");
          formCheck.className = "form-check";

          const checkbox = document.createElement("input");
          checkbox.type = "checkbox";
          checkbox.className =
            "form-check-input member-checkbox cursor-pointer";
          checkbox.id = `member_check_${member.id}`;
          checkbox.value = member.id;
          checkbox.dataset.name = member.name;

          const label = document.createElement("label");
          label.className = "form-check-label w-100 cursor-pointer ms-1";
          label.htmlFor = `member_check_${member.id}`;
          label.textContent = member.name;

          formCheck.appendChild(checkbox);
          formCheck.appendChild(label);
          div.appendChild(formCheck);
          memberCheckboxList.appendChild(div);

          checkbox.addEventListener("change", function () {
            if (this.checked) {
              selectedMembers.set(this.value, this.dataset.name);
            } else {
              selectedMembers.delete(this.value);
            }
            updateSelectedMembersUI();
          });
        });
        memberList.classList.remove("d-none");
        memberPlaceholder.classList.add("d-none");
      } else {
        memberList.classList.add("d-none");
        memberPlaceholder.textContent = "No active members found in this team.";
        memberPlaceholder.classList.remove("d-none");
      }
    } catch (error) {
      console.error("Error loading team members:", error);
      Toast.error("Failed to load team members");
    }
  }

  function updateSelectedMembersUI() {
    // Clear container except for the placeholder
    selectedMembersContainer.innerHTML = "";

    if (selectedMembers.size === 0) {
      selectedMembersContainer.appendChild(noMembersSelected);
      noMembersSelected.classList.remove("d-none");
      return;
    }

    noMembersSelected.classList.add("d-none");

    selectedMembers.forEach((name, id) => {
      const item = document.createElement("div");
      item.className =
        "list-group-item d-flex justify-content-between align-items-center animate__animated animate__fadeIn";
      item.innerHTML = `
        <span>${name}</span>
        <button type="button" class="btn btn-sm btn-outline-danger remove-member-btn" data-id="${id}">
          <i class="bi bi-trash"></i>
        </button>
      `;

      item
        .querySelector(".remove-member-btn")
        .addEventListener("click", function () {
          selectedMembers.delete(id);
          // Uncheck the checkbox in the dropdown
          const checkbox = document.getElementById(`member_check_${id}`);
          if (checkbox) {
            checkbox.checked = false;
          }
          updateSelectedMembersUI();
        });

      selectedMembersContainer.appendChild(item);
    });
  }

  // Handle form submission
  createProjectForm.addEventListener("submit", async function (e) {
    e.preventDefault();

    const formData = new FormData(createProjectForm);
    const data = {
      name: formData.get("name"),
      abbreviation: formData.get("abbreviation"),
      start_date: formData.get("start_date") || null,
      end_date: formData.get("end_date") || null,
      leader_id: parseInt(leaderID.value),
      team_id: parseInt(formData.get("team_id")),
      member_ids: Array.from(selectedMembers.keys()).map((id) => parseInt(id)),
    };

    // Basic validation
    if (!data.name || !data.abbreviation || !data.leader_id || !data.team_id) {
      Toast.error("Please fill in all required fields");
      return;
    }

    createProjectBtn.disabled = true;
    createProjectBtn.innerHTML =
      '<span class="spinner-border spinner-border-sm me-1"></span>Creating...';

    try {
      const response = await AdminProjectService.createProject(data);
      Toast.success(response.message || "Project created successfully");
      setTimeout(() => {
        window.location.href = "/admin/projects";
      }, 1500);
    } catch (error) {
      console.error("Error creating project:", error);
      let errorMessage = error.message || "Failed to create project";
      if (error.details) {
        const details = Object.entries(error.details)
          .map(([field, msg]) => `${field} ${msg}`)
          .join("<br>");
        errorMessage = `<strong>${errorMessage}</strong>:<br>${details}`;
      }
      Toast.error(errorMessage);
      createProjectBtn.disabled = false;
      createProjectBtn.innerHTML =
        '<i class="bi bi-check-circle me-1"></i>Create Project';
    }
  });
});
