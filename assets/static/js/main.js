document.addEventListener('DOMContentLoaded', function() {
    const sidebar = document.getElementById('sidebar');
    const content = document.getElementById('content');
    const sidebarToggle = document.getElementById('sidebarToggle');

    // Load sidebar state from localStorage
    const sidebarState = localStorage.getItem('sidebarState');
    if (sidebarState === 'active') {
        sidebar.classList.add('active');
        content.classList.add('active');
    }

    sidebarToggle.addEventListener('click', function() {
        // Apply transition to both elements
        sidebar.style.transition = 'all 0.3s';
        content.style.transition = 'all 0.3s';

        // Toggle the active class
        sidebar.classList.toggle('active');
        content.classList.toggle('active');

        // Save or remove sidebar state to/from localStorage
        if (sidebar.classList.contains('active')) {
            localStorage.setItem('sidebarState', 'active');
        } else {
            localStorage.removeItem('sidebarState');
        }

        // Remove transition property after animation ends
        setTimeout(() => {
            sidebar.style.transition = '';
            content.style.transition = '';
        }, 300); // Duration of the transition in milliseconds
    });
});
