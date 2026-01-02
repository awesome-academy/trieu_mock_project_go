document.addEventListener("DOMContentLoaded", function () {
  const skillListContainer = document.getElementById("skillListContainer");
  const loadingTemplate = document.getElementById("loadingTemplate");

  async function loadSkills(offset = 0) {
    const limit = 10;

    // Show loading spinner
    skillListContainer.innerHTML = loadingTemplate.innerHTML;

    try {
      const html = await AdminSkillService.searchSkills({
        limit,
        offset,
      });
      skillListContainer.innerHTML = html;
      attachEvents();
    } catch (error) {
      console.error("Error loading skills:", error);
      Toast.error("Failed to load skills list");
      skillListContainer.innerHTML =
        '<div class="alert alert-danger">Failed to load skills.</div>';
    }
  }

  function attachEvents() {
    // Pagination events
    const paginationLinks = skillListContainer.querySelectorAll(".page-link");
    paginationLinks.forEach((link) => {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        const offsetAttr = this.getAttribute("data-offset");
        if (offsetAttr !== null) {
          const offset = parseInt(offsetAttr, 10);
          if (!Number.isNaN(offset) && offset >= 0) {
            loadSkills(offset);
          }
        }
      });
    });

    // Delete events
    const deleteBtns = skillListContainer.querySelectorAll(".delete-skill-btn");
    deleteBtns.forEach((btn) => {
      btn.addEventListener("click", async function () {
        const id = this.getAttribute("data-id");
        const name = this.getAttribute("data-name");

        if (
          confirm(
            `Are you sure you want to delete skill "${escapeForDialog(name)}"?`
          )
        ) {
          try {
            const response = await AdminSkillService.deleteSkill(id);
            Toast.success(response.message || "Skill deleted successfully");
            loadSkills(0);
          } catch (error) {
            console.error("Error deleting skill:", error);
            Toast.error(error.message || "Failed to delete skill");
          }
        }
      });
    });
  }

  // Initial load
  loadSkills(0);
});

function escapeForDialog(str) {
  return str
    .replace(/\\/g, "\\\\")
    .replace(/"/g, '\\"')
    .replace(/\n/g, "\\n")
    .replace(/\r/g, "\\r");
}
