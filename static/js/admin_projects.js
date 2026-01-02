document.addEventListener("DOMContentLoaded", function () {
  const projectListContainer = document.getElementById("projectListContainer");
  const loadingTemplate = document.getElementById("loadingTemplate");
  const teamFilter = document.getElementById("teamFilter");

  async function loadProjects(offset = 0) {
    projectListContainer.innerHTML = loadingTemplate.innerHTML;
    try {
      const params = {
        limit: 10,
        offset: offset,
      };

      if (teamFilter && teamFilter.value) {
        params.team_id = teamFilter.value;
      }

      const html = await AdminProjectService.listProjects(params);
      projectListContainer.innerHTML = html;
      attachEventListeners();
    } catch (error) {
      console.error("Error loading projects:", error);
      Toast.error("Failed to load projects list");
      projectListContainer.innerHTML = `<div class="alert alert-danger">Failed to load projects.</div>`;
    }
  }

  if (teamFilter) {
    teamFilter.addEventListener("change", function () {
      loadProjects(0);
    });
  }

  function attachEventListeners() {
    // Pagination
    document.querySelectorAll(".page-link").forEach((link) => {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        const offset = this.getAttribute("data-offset");
        if (offset !== null) {
          loadProjects(parseInt(offset));
        }
      });
    });

    // Delete buttons
    document.querySelectorAll(".delete-project-btn").forEach((btn) => {
      btn.addEventListener("click", async function () {
        const id = this.getAttribute("data-id");
        const name = this.getAttribute("data-name");

        if (confirm(`Are you sure you want to delete project "${name}"?`)) {
          try {
            const response = await AdminProjectService.deleteProject(id);
            Toast.success(response.message || "Project deleted successfully");
            loadProjects();
          } catch (error) {
            console.error("Error deleting project:", error);
            Toast.error(error.message || "Failed to delete project");
          }
        }
      });
    });
  }

  loadProjects();
});
