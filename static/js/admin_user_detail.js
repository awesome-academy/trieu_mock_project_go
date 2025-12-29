document.addEventListener("DOMContentLoaded", function () {
  const deleteUserBtn = document.getElementById("deleteUserBtn");

  if (deleteUserBtn) {
    deleteUserBtn.addEventListener("click", async function () {
      const userId = this.dataset.userId;
      const userName = this.dataset.userName;

      if (
        confirm(
          `Are you sure you want to delete user "${userName}"? This action cannot be undone.`
        )
      ) {
        try {
          const response = await fetch(`/admin/users/${userId}`, {
            method: "DELETE",
            headers: {
              "Content-Type": "application/json",
            },
          });

          if (response.ok) {
            window.location.href = "/admin/users";
          } else {
            const err = await response.json();
            alert("Error: " + (err.error || "Failed to delete user"));
          }
        } catch (error) {
          console.error("Error deleting user:", error);
          alert("An error occurred while deleting the user.");
        }
      }
    });
  }
});
