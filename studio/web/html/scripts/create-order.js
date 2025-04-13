document.addEventListener('DOMContentLoaded', function () {
  const checkboxes = document.querySelectorAll('input[name="model_ids"]');
  const submitButton = document.getElementById('submit-button');

  function toggleSubmitButton() {
    const isChecked = Array.from(checkboxes).some(checkbox => checkbox.checked);
    submitButton.disabled = !isChecked;
  }

  checkboxes.forEach(checkbox => {
    checkbox.addEventListener('change', toggleSubmitButton);
  });

  // Проверим сразу при загрузке страницы
  toggleSubmitButton();
});
