document.addEventListener('DOMContentLoaded', function () {
  const modal = document.getElementById('modal');
  const modalBody = document.getElementById('modal-body');
  const modalHeader = document.getElementById('modal-header');
  const header = modalHeader.innerText
  const closeModal = document.querySelector('#modal .close');

  function openModal(content = 'Загрузка...') {
    modalBody.innerHTML = content;
    modal.style.display = 'flex';
  }

  function closeModalWindow() {
    modal.style.display = 'none';
  }

  // Закрытие модалки
  closeModal.addEventListener('click', closeModalWindow);
  window.addEventListener('click', function (e) {
    if (e.target === modal) closeModalWindow();
  });
  window.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') closeModalWindow();
  });

  // Открытие модалки по клику на ссылку модели
  document.querySelectorAll('.item-tr').forEach(link => {
    link.addEventListener('click', function (e) {
      e.preventDefault();
      const id = this.dataset.id;
      const url = this.dataset.url;
      modalHeader.innerText = `${header}${id}`
      openModal();

      // Fetch данных о составе модели
      fetch(`/${url}?id=${id}`)
        .then(response => response.ok ? response.text() : Promise.reject('Ошибка загрузки'))
        .then(html => modalBody.innerHTML = html)
        .catch(error => modalBody.innerHTML = `<p style="color: red;">${error}</p>`);
    });
  });
});