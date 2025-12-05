// ===================== SHARED SIDEBAR & HEADER FUNCTIONALITY =====================

$(document).ready(function() {
    let sidebarExpanded = false;

    // Toggle Sidebar Function
    function toggleSidebar() {
        sidebarExpanded = !sidebarExpanded;
        $('#sidebar').toggleClass('expanded', sidebarExpanded);
        $('#content').toggleClass('expanded', sidebarExpanded);

        if ($(window).width() <= 768) {
            $('#sidebar').toggleClass('mobile-active', sidebarExpanded);
            $('.mobile-overlay').toggleClass('active', sidebarExpanded);
        }
    }

    // Sidebar Toggle Button Click
    $('#sidebarToggle').on('click', function(e) {
        e.stopPropagation();
        toggleSidebar();
    });

    // Close Sidebar on Document Click (Mobile)
    $(document).on('click', function(e) {
        if ($(window).width() <= 768 && sidebarExpanded) {
            if (!$(e.target).closest('#sidebar').length && !$(e.target).is('#sidebarToggle') && !$(e.target).closest('#sidebarToggle').length) {
                toggleSidebar();
            }
        }
    });

    // Mobile Overlay Click
    $('.mobile-overlay').on('click', function() {
        if (sidebarExpanded) {
            toggleSidebar();
        }
    });

    // Prevent Sidebar from Closing when Clicking Inside
    $('#sidebar').on('click', function(e) {
        e.stopPropagation();
    });

    // Check Mobile View on Resize
    function checkMobileView() {
        if($(window).width() <= 768) {
            $('#sidebar').removeClass('expanded').removeClass('mobile-active');
            $('#content').removeClass('expanded');
            $('.mobile-overlay').removeClass('active');
            sidebarExpanded = false;
        }
    }

    checkMobileView();
    $(window).resize(checkMobileView);

    // --- GENEL SİSTEM İŞLEVLERİ BAŞLANGICI ---
    
    // ===================== TÜM SUBMIT BUTONLARINI YÖNETEN GENEL İŞLEV (TEK KURAL) =====================
    
    const allForms = document.querySelectorAll('form');
    // Font Awesome yükleme ikonu ve "Bekleyiniz..." metni
    const loadingHtml = '<i class="fas fa-spinner fa-spin me-2"></i>Bekleyiniz...'; 

    allForms.forEach(form => {
        // Formun submit olayını dinle
        $(form).on('submit', function(e) {
            // Formun içindeki submit tipindeki ilk düğmeyi bul (button veya input)
            const submitButton = $(this).find('button[type="submit"], input[type="submit"]').first();

            if (submitButton.length) {
                // Buton zaten devre dışıysa, işlemi tekrar çalıştırma (çift tıklamayı engeller)
                if (submitButton.prop('disabled')) {
                    // Buton devre dışıysa, formu submit etmeyi durdur (çift gönderimi engellemenin son adımı)
                    e.preventDefault(); 
                    return;
                }

                // Orijinal içeriği kaydet (sayfa yenilenmese bile, geri yükleme ihtimali için tutulur)
                if (!submitButton.data('original-content')) {
                    submitButton.data('original-content', submitButton.html() || submitButton.val());
                }

                // Butonu devre dışı bırak ve metni değiştir
                submitButton.prop('disabled', true);
                
                if (submitButton.is('button')) {
                    submitButton.html(loadingHtml);
                } else if (submitButton.is('input')) {
                    submitButton.val('Bekleyiniz...');
                }
            }
            // NOT: e.preventDefault() burada çağrılmaz. Butonu devre dışı bırakıp metni değiştirdikten sonra,
            // form normal submit akışına (sunucuya gönderim ve sayfa yenileme) devam eder.
        });
    });
    // ======================================================================================================

    // ===================== KALDIRILAN AUTH VE USER FORM İŞLEVLERİ =====================
    // Auth Form Submit kaldırıldı.
    // User Form Submit kaldırıldı.
    // validatePassword kaldırıldı.
    // Bu formlar artık sadece genel submit kuralına uyar ve submit olduklarında
    // direkt backend'e gönderilirler.

    // Input Icon Focus Effects
    $('.form-control').on('focus', function() {
        $(this).prev('.input-icon').css({
            'border-color': 'var(--primary-color)',
            'background-color': '#eef2ff',
            'color': 'var(--primary-color)'
        });
    });

    $('.form-control').on('blur', function() {
        $(this).prev('.input-icon').css({
            'border-color': '#e5e7eb',
            'background-color': '#f1f5f9',
            'color': '#64748b'
        });
    });

    // ===================== PANEL CHARTS =====================
    
    // Sales Chart (Panel Page)
    if ($('#salesChart').length) {
        const salesCtx = document.getElementById('salesChart').getContext('2d');
        const salesChart = new Chart(salesCtx, {
            type: 'line',
            data: {
                labels: ['Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran'],
                datasets: [{
                    label: 'Satışlar',
                    data: [12000, 19000, 15000, 25000, 22000, 30000],
                    borderColor: '#6366f1',
                    backgroundColor: 'rgba(99, 102, 241, 0.1)',
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        grid: {
                            drawBorder: false
                        }
                    },
                    x: {
                        grid: {
                            display: false
                        }
                    }
                }
            }
        });
    }

    // Category Chart (Panel Page)
    if ($('#categoryChart').length) {
        const categoryCtx = document.getElementById('categoryChart').getContext('2d');
        const categoryChart = new Chart(categoryCtx, {
            type: 'doughnut',
            data: {
                labels: ['Elektronik', 'Giyim', 'Ev Eşyaları', 'Kitaplar', 'Diğer'],
                datasets: [{
                    data: [30, 25, 20, 15, 10],
                    backgroundColor: [
                        '#6366f1',
                        '#8b5cf6',
                        '#10b981',
                        '#f59e0b',
                        '#06b6d4'
                    ],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    }

});