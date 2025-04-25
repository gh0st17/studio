document.addEventListener('DOMContentLoaded', () => {
  const table = document.getElementById('sort-table');
  const headers = table.querySelectorAll('th');
  const tbody = table.querySelector('tbody');

  window.addEventListener('DOMContentLoaded', () => {
    const index = parseInt(localStorage.getItem('index'));
    const ascending = localStorage.getItem('ascending') == 'true';
  
    if (!isNaN(index) && !isNaN(ascending)) {
      sortTableByColumn(index, ascending);
    } else {
      sortTableByColumn(0, true);
    }
  });

  headers.forEach(th => {
    th.dataset.label = th.textContent.trim();
    th.dataset.order = '';
  });

  function updateHeaderArrows(activeIndex, ascending) {
    headers.forEach((th, i) => {
      th.textContent = th.dataset.label;
  
      let arrow = th.querySelector('.sort-arrow');
      if (!arrow) {
        arrow = document.createElement('span');
        arrow.classList.add('sort-arrow');
        th.appendChild(arrow);
      }
  
      if (i === activeIndex) {
        arrow.textContent = ascending ? '▲' : '▼';
        th.dataset.order = ascending ? 'asc' : 'desc';
      } else {
        arrow.textContent = ' ';
        th.dataset.order = '';
      }
    });
  }  

  function parseValue(val) {
    if (val === '---' || !val) return null;

    const num = parseFloat(val.replace(/[^\d.-]/g, ''));
    if (!isNaN(num)) return num;

    const date = Date.parse(val);
    if (!isNaN(date)) return date;

    return val.toLowerCase();
  }

  function sortTableByColumn(index, ascending) {
    const rows = Array.from(tbody.querySelectorAll('tr'));
    rows.sort((a, b) => {
      const aCell = a.children[index].textContent.trim();
      const bCell = b.children[index].textContent.trim();

      const aVal = parseValue(aCell);
      const bVal = parseValue(bCell);

      return ascending
        ? aVal > bVal ? 1 : aVal < bVal ? -1 : 0
        : aVal < bVal ? 1 : aVal > bVal ? -1 : 0;
    });

    rows.forEach(row => tbody.appendChild(row));
    localStorage.setItem('index', index);
    localStorage.setItem('ascending', ascending);
    updateHeaderArrows(index, ascending);
  }

  // Обработчики кликов по заголовкам
  headers.forEach((th, index) => {
    th.addEventListener('click', () => {
      const ascending = th.dataset.order !== 'asc';
      sortTableByColumn(index, ascending);
    });
  });
});
