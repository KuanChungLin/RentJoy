let navItems = document.querySelectorAll('.nav-item');

window.addEventListener('load',function () {
    navItems.forEach(navItem => {
        navItem.onclick = function (event) {
            let screenWidth = window.innerWidth;
                
            if (screenWidth < 1025 && this.querySelector('.sub-menu')) {
                let color = getComputedStyle(navItem).backgroundColor;
                navItem.style.backgroundColor = color === "rgba(0, 0, 0, 0)" ? "#f6f6f6" : "rgba(0, 0, 0, 0)";
            }


            let allSubMenus = document.querySelectorAll('.sub-menu');
            allSubMenus.forEach(subMenu => {
                if (subMenu !== this.querySelector('.sub-menu')) {
                    subMenu.style.display = 'none';
                }
            });


            let currentSubMenu = this.querySelector('.sub-menu');
            if (currentSubMenu) {
                currentSubMenu.style.display = currentSubMenu.style.display === 'none' ? 'block' : 'none';
            }
            

            event.stopPropagation();
        };
    });
    

    
    // 當觸發.faqBtn:hover時，關閉所有.sub-menu
    let faqBtn = document.querySelector('.faq-btn');
    faqBtn.onmouseover = function () {
        if (window.innerWidth > 1024) {
            let subMenus = document.querySelectorAll('.sub-menu');
            subMenus.forEach(subMenu => subMenu.style.display = 'none');
        }
    };

    // 點擊在.nav-item以外的地方，則隱藏所有.sub-menu
    document.onclick = function () {
        if (window.innerWidth > 1024) {
            let subMenus = document.querySelectorAll('.sub-menu');
            subMenus.forEach(subMenu => subMenu.style.display = 'none');
        }
    }

    // 視窗大小變動時，隱藏所有 .sub-menu
    window.onresize = function () {
        if (window.innerWidth > 1025) {
            let subMenus = document.querySelectorAll('.sub-menu');
            subMenus.forEach(subMenu => subMenu.style.display = 'none');
        }
    }
});