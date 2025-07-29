/**
 * QWQ Blog Template - Main JavaScript
 * Author: xiaolin
 */

document.addEventListener('DOMContentLoaded', function() {
    // 初始化导航菜单
    initNavigation();
    
    // 初始化搜索模态框
    initSearchModal();
    
    // 初始化FAQ折叠面板（如果存在）
    initFaqAccordion();
    
    // 初始化文章目录（如果存在）
    initTableOfContents();
    
    // 初始化代码高亮（如果存在）
    initCodeHighlight();
    
    // 初始化暗色模式切换
    initDarkModeToggle();
    
    // 初始化滚动到顶部按钮
    initScrollToTop();
    
    // 初始化图片懒加载
    initLazyLoading();
    
    // 初始化评论功能（如果存在）
    initComments();
    
    // 初始化Popular区域交互功能
    initPopularSection();
    
    // 初始化页面动画效果
    initPageAnimations();
});

/**
 * 初始化导航菜单
 */
function initNavigation() {
    const navToggle = document.querySelector('.nav-toggle');
    const mainNav = document.querySelector('.main-nav');
    
    if (navToggle && mainNav) {
        navToggle.addEventListener('click', function() {
            mainNav.classList.toggle('active');
            document.body.classList.toggle('nav-open');
        });
        
        // 点击导航链接后关闭移动导航菜单
        const navLinks = mainNav.querySelectorAll('a');
        navLinks.forEach(link => {
            link.addEventListener('click', function() {
                if (window.innerWidth < 768) {
                    mainNav.classList.remove('active');
                    document.body.classList.remove('nav-open');
                }
            });
        });
        
        // 点击页面其他区域关闭导航菜单
        document.addEventListener('click', function(event) {
            if (window.innerWidth < 768 && 
                !event.target.closest('.main-nav') && 
                !event.target.closest('.nav-toggle') && 
                mainNav.classList.contains('active')) {
                mainNav.classList.remove('active');
                document.body.classList.remove('nav-open');
            }
        });
    }
    
    // 设置当前页面的导航链接为活动状态
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.main-nav a');
    
    navLinks.forEach(link => {
        const linkPath = link.getAttribute('href');
        if (currentPath === linkPath || 
            (linkPath !== '/' && currentPath.startsWith(linkPath))) {
            link.classList.add('active');
        }
    });
}

/**
 * 初始化搜索模态框
 */
function initSearchModal() {
    const searchBtn = document.querySelector('.search-btn');
    const searchModal = document.querySelector('.search-modal');
    const closeSearch = document.querySelector('.close-search');
    const searchForm = document.querySelector('.search-form');
    
    if (searchBtn && searchModal) {
        searchBtn.addEventListener('click', function() {
            searchModal.classList.add('active');
            document.body.style.overflow = 'hidden';
            setTimeout(() => {
                searchModal.querySelector('input').focus();
            }, 100);
        });
        
        if (closeSearch) {
            closeSearch.addEventListener('click', function() {
                searchModal.classList.remove('active');
                document.body.style.overflow = '';
            });
        }
        
        // 点击模态框外部关闭
        searchModal.addEventListener('click', function(event) {
            if (event.target === searchModal) {
                searchModal.classList.remove('active');
                document.body.style.overflow = '';
            }
        });
        
        // ESC键关闭搜索
        document.addEventListener('keydown', function(event) {
            if (event.key === 'Escape' && searchModal.classList.contains('active')) {
                searchModal.classList.remove('active');
                document.body.style.overflow = '';
            }
        });
        
        // 搜索表单提交
        if (searchForm) {
            searchForm.addEventListener('submit', function(event) {
                event.preventDefault();
                const searchInput = searchForm.querySelector('input').value.trim();
                if (searchInput) {
                    // 这里可以实现搜索逻辑，例如跳转到搜索结果页面
                    window.location.href = `/search?q=${encodeURIComponent(searchInput)}`;
                }
            });
        }
    }
}

/**
 * 初始化FAQ折叠面板
 */
function initFaqAccordion() {
    const faqQuestions = document.querySelectorAll('.faq-question');
    
    if (faqQuestions.length > 0) {
        faqQuestions.forEach(question => {
            question.addEventListener('click', function() {
                const answer = this.nextElementSibling;
                const isActive = this.classList.contains('active');
                
                // 关闭所有其他FAQ
                faqQuestions.forEach(q => {
                    if (q !== this) {
                        q.classList.remove('active');
                        q.nextElementSibling.classList.remove('active');
                    }
                });
                
                // 切换当前FAQ状态
                if (isActive) {
                    this.classList.remove('active');
                    answer.classList.remove('active');
                } else {
                    this.classList.add('active');
                    answer.classList.add('active');
                }
            });
        });
    }
}

/**
 * 初始化文章目录
 */
function initTableOfContents() {
    const articleContent = document.querySelector('.article-text');
    const tableOfContents = document.querySelector('.toc');
    
    if (articleContent && tableOfContents) {
        const headings = articleContent.querySelectorAll('h2, h3, h4');
        if (headings.length > 0) {
            // 清空现有目录
            tableOfContents.innerHTML = '';
            
            // 创建目录项
            const tocList = document.createElement('ul');
            let currentLevel2List = null;
            let currentLevel3List = null;
            
            headings.forEach((heading, index) => {
                const headingText = heading.textContent;
                const headingId = `heading-${index}`;
                heading.setAttribute('id', headingId);
                
                const listItem = document.createElement('li');
                const link = document.createElement('a');
                link.href = `#${headingId}`;
                link.textContent = headingText;
                
                // 平滑滚动到目标位置
                link.addEventListener('click', function(e) {
                    e.preventDefault();
                    const targetHeading = document.getElementById(headingId);
                    const headerOffset = 100; // 考虑固定头部的高度
                    const elementPosition = targetHeading.getBoundingClientRect().top;
                    const offsetPosition = elementPosition + window.pageYOffset - headerOffset;
                    
                    window.scrollTo({
                        top: offsetPosition,
                        behavior: 'smooth'
                    });
                });
                
                if (heading.tagName === 'H2') {
                    listItem.appendChild(link);
                    tocList.appendChild(listItem);
                    currentLevel2List = document.createElement('ul');
                    listItem.appendChild(currentLevel2List);
                    currentLevel3List = null;
                } else if (heading.tagName === 'H3') {
                    if (!currentLevel2List) {
                        currentLevel2List = document.createElement('ul');
                        tocList.appendChild(currentLevel2List);
                    }
                    listItem.appendChild(link);
                    currentLevel2List.appendChild(listItem);
                    currentLevel3List = document.createElement('ul');
                    listItem.appendChild(currentLevel3List);
                } else if (heading.tagName === 'H4') {
                    if (!currentLevel3List) {
                        if (!currentLevel2List) {
                            currentLevel2List = document.createElement('ul');
                            tocList.appendChild(currentLevel2List);
                        }
                        const parentItem = document.createElement('li');
                        currentLevel2List.appendChild(parentItem);
                        currentLevel3List = document.createElement('ul');
                        parentItem.appendChild(currentLevel3List);
                    }
                    listItem.appendChild(link);
                    currentLevel3List.appendChild(listItem);
                }
            });
            
            tableOfContents.appendChild(tocList);
            
            // 监听滚动，高亮当前目录项
            window.addEventListener('scroll', highlightTableOfContents);
        } else {
            // 如果没有标题，隐藏目录小部件
            const tocWidget = tableOfContents.closest('.widget');
            if (tocWidget) {
                tocWidget.style.display = 'none';
            }
        }
    }
}

/**
 * 高亮当前可见的目录项
 */
function highlightTableOfContents() {
    const headings = document.querySelectorAll('.article-text h2, .article-text h3, .article-text h4');
    const tocLinks = document.querySelectorAll('.toc a');
    
    if (headings.length === 0 || tocLinks.length === 0) return;
    
    // 找到当前可见的标题
    let currentHeadingIndex = -1;
    const scrollPosition = window.scrollY;
    const headerOffset = 120; // 考虑固定头部的高度和一些额外空间
    
    for (let i = 0; i < headings.length; i++) {
        const headingTop = headings[i].getBoundingClientRect().top + window.pageYOffset - headerOffset;
        if (scrollPosition >= headingTop) {
            currentHeadingIndex = i;
        } else {
            break;
        }
    }
    
    // 移除所有高亮
    tocLinks.forEach(link => {
        link.classList.remove('active');
    });
    
    // 高亮当前项
    if (currentHeadingIndex >= 0) {
        const currentHeadingId = headings[currentHeadingIndex].getAttribute('id');
        const currentTocLink = document.querySelector(`.toc a[href="#${currentHeadingId}"]`);
        if (currentTocLink) {
            currentTocLink.classList.add('active');
        }
    }
}

/**
 * 初始化代码高亮
 */
function initCodeHighlight() {
    // 检查是否已加载Prism.js
    if (typeof Prism !== 'undefined') {
        Prism.highlightAll();
    }
}

/**
 * 初始化暗色模式切换
 */
function initDarkModeToggle() {
    const themeSwitch = document.querySelector('.theme-switch');
    
    if (themeSwitch) {
        // 检查本地存储中的主题偏好
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === 'dark') {
            document.body.classList.add('dark-mode');
            updateThemeIcon(true);
        }
        
        themeSwitch.addEventListener('click', function() {
            const isDarkMode = document.body.classList.toggle('dark-mode');
            localStorage.setItem('theme', isDarkMode ? 'dark' : 'light');
            updateThemeIcon(isDarkMode);
        });
    }
}

/**
 * 更新主题图标
 */
function updateThemeIcon(isDarkMode) {
    const themeIcon = document.querySelector('.theme-switch i');
    if (themeIcon) {
        if (isDarkMode) {
            themeIcon.className = 'fas fa-sun';
        } else {
            themeIcon.className = 'fas fa-moon';
        }
    }
}

/**
 * 初始化滚动到顶部按钮
 */
function initScrollToTop() {
    // 创建滚动到顶部按钮
    const scrollTopBtn = document.createElement('button');
    scrollTopBtn.className = 'scroll-top-btn';
    scrollTopBtn.innerHTML = '<i class="fas fa-arrow-up"></i>';
    document.body.appendChild(scrollTopBtn);
    
    // 监听滚动事件，显示/隐藏按钮
    window.addEventListener('scroll', function() {
        if (window.pageYOffset > 300) {
            scrollTopBtn.classList.add('active');
        } else {
            scrollTopBtn.classList.remove('active');
        }
    });
    
    // 点击按钮滚动到顶部
    scrollTopBtn.addEventListener('click', function() {
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    });
    
    // 添加按钮样式
    const style = document.createElement('style');
    style.textContent = `
        .scroll-top-btn {
            position: fixed;
            bottom: 30px;
            right: 30px;
            width: 40px;
            height: 40px;
            border-radius: 50%;
            background-color: var(--primary-color);
            color: white;
            border: none;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            opacity: 0;
            visibility: hidden;
            transition: opacity 0.3s, visibility 0.3s, background-color 0.3s;
            z-index: 99;
        }
        
        .scroll-top-btn.active {
            opacity: 1;
            visibility: visible;
        }
        
        .scroll-top-btn:hover {
            background-color: var(--primary-dark);
        }
        
        @media (max-width: 768px) {
            .scroll-top-btn {
                bottom: 20px;
                right: 20px;
                width: 36px;
                height: 36px;
            }
        }
    `;
    document.head.appendChild(style);
}

/**
 * 初始化图片懒加载
 */
function initLazyLoading() {
    // 检查浏览器是否支持IntersectionObserver
    if ('IntersectionObserver' in window) {
        const lazyImages = document.querySelectorAll('img[data-src]');
        
        const imageObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    img.src = img.dataset.src;
                    img.removeAttribute('data-src');
                    imageObserver.unobserve(img);
                }
            });
        });
        
        lazyImages.forEach(img => {
            imageObserver.observe(img);
        });
    } else {
        // 回退方案：简单的滚动事件监听
        let lazyImages = document.querySelectorAll('img[data-src]');
        let active = false;
        
        const lazyLoad = function() {
            if (active === false) {
                active = true;
                
                setTimeout(function() {
                    lazyImages.forEach(function(lazyImage) {
                        if ((lazyImage.getBoundingClientRect().top <= window.innerHeight && lazyImage.getBoundingClientRect().bottom >= 0) && getComputedStyle(lazyImage).display !== 'none') {
                            lazyImage.src = lazyImage.dataset.src;
                            lazyImage.removeAttribute('data-src');
                            
                            lazyImages = Array.from(lazyImages).filter(function(image) {
                                return image !== lazyImage;
                            });
                            
                            if (lazyImages.length === 0) {
                                document.removeEventListener('scroll', lazyLoad);
                                window.removeEventListener('resize', lazyLoad);
                                window.removeEventListener('orientationchange', lazyLoad);
                            }
                        }
                    });
                    
                    active = false;
                }, 200);
            }
        };
        
        document.addEventListener('scroll', lazyLoad);
        window.addEventListener('resize', lazyLoad);
        window.addEventListener('orientationchange', lazyLoad);
        lazyLoad();
    }
}

/**
 * 初始化评论功能
 */
function initComments() {
    const commentForm = document.querySelector('.comment-form');
    const commentsList = document.querySelector('.comments-list');
    
    if (commentForm && commentsList) {
        commentForm.addEventListener('submit', function(event) {
            event.preventDefault();
            
            const commentTextarea = commentForm.querySelector('textarea');
            const commentText = commentTextarea.value.trim();
            
            if (commentText) {
                // 在实际应用中，这里应该发送AJAX请求到服务器
                // 这里仅作为前端演示
                addComment({
                    avatar: '/images/user-avatar.jpg', // 假设用户头像
                    author: '当前用户',
                    date: new Date().toLocaleDateString(),
                    text: commentText
                });
                
                // 清空评论框
                commentTextarea.value = '';
            }
        });
        
        // 回复按钮点击事件
        commentsList.addEventListener('click', function(event) {
            if (event.target.classList.contains('reply-btn') || 
                event.target.parentElement.classList.contains('reply-btn')) {
                
                const commentEl = event.target.closest('.comment');
                const authorName = commentEl.querySelector('.comment-author').textContent;
                
                const commentTextarea = commentForm.querySelector('textarea');
                commentTextarea.value = `@${authorName} `;
                commentTextarea.focus();
                
                // 滚动到评论表单
                commentForm.scrollIntoView({ behavior: 'smooth' });
            }
        });
    }
}

/**
 * 添加评论到评论列表
 */
function addComment(comment) {
    const commentsList = document.querySelector('.comments-list');
    
    if (commentsList) {
        const commentEl = document.createElement('div');
        commentEl.className = 'comment';
        
        commentEl.innerHTML = `
            <div class="comment-avatar">
                <img src="${comment.avatar}" alt="${comment.author}">
            </div>
            <div class="comment-content">
                <div class="comment-header">
                    <h4 class="comment-author">${comment.author}</h4>
                    <span class="comment-date">${comment.date}</span>
                </div>
                <div class="comment-text">
                    <p>${comment.text}</p>
                </div>
                <div class="comment-actions">
                    <a href="#" class="reply-btn"><i class="fas fa-reply"></i> 回复</a>
                    <span class="like-btn"><i class="far fa-heart"></i> 0</span>
                </div>
            </div>
        `;
        
        // 添加点赞功能
        const likeBtn = commentEl.querySelector('.like-btn');
        likeBtn.addEventListener('click', function() {
            const likeCount = this.textContent.trim().split(' ')[1];
            let newCount = parseInt(likeCount) + 1;
            this.innerHTML = `<i class="fas fa-heart"></i> ${newCount}`;
            this.style.color = 'var(--error-color)';
            this.style.pointerEvents = 'none';
        });
        
        // 将新评论添加到列表顶部
        commentsList.insertBefore(commentEl, commentsList.firstChild);
    }
}

/**
 * 初始化Popular区域交互功能
 */
function initPopularSection() {
    const popularCards = document.querySelectorAll('.popular-card');
    
    if (popularCards.length > 0) {
        popularCards.forEach(card => {
            // 添加点击事件
            card.addEventListener('click', function() {
                // 这里可以添加跳转到文章详情页的逻辑
                const title = this.querySelector('.popular-title').textContent;
                console.log('点击了文章:', title);
                // 实际应用中可以跳转到文章页面
                // window.location.href = '/article/' + articleId;
            });
            
            // 添加键盘导航支持
            card.setAttribute('tabindex', '0');
            card.addEventListener('keydown', function(e) {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    this.click();
                }
            });
            
            // 添加鼠标进入和离开的音效反馈（可选）
            card.addEventListener('mouseenter', function() {
                this.style.cursor = 'pointer';
            });
        });
        
        // 添加Popular区域的加载动画
        const popularSection = document.querySelector('.popular');
        if (popularSection) {
            // 使用Intersection Observer来触发动画
            const observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        entry.target.classList.add('animate-in');
                        
                        // 为每个卡片添加延迟动画
                        const cards = entry.target.querySelectorAll('.popular-card');
                        cards.forEach((card, index) => {
                            setTimeout(() => {
                                card.classList.add('animate-card');
                            }, index * 100);
                        });
                        
                        observer.unobserve(entry.target);
                    }
                });
            }, {
                threshold: 0.1
            });
            
            observer.observe(popularSection);
        }
    }
}

/**
 * 初始化页面动画效果
 */
function initPageAnimations() {
    // 为所有区域添加滚动动画
    const animatedSections = document.querySelectorAll('.features, .newest, .popular, .cta');
    
    if (animatedSections.length > 0 && 'IntersectionObserver' in window) {
        const sectionObserver = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('section-animate');
                    sectionObserver.unobserve(entry.target);
                }
            });
        }, {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        });
        
        animatedSections.forEach(section => {
            sectionObserver.observe(section);
        });
    }
    
    // 为文章卡片添加悬停效果增强
    const allCards = document.querySelectorAll('.article-card, .feature-card, .popular-card');
    allCards.forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-8px) scale(1.02)';
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = '';
        });
    });
    
    // 添加平滑滚动到锚点
    const anchorLinks = document.querySelectorAll('a[href^="#"]');
    anchorLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            const targetId = this.getAttribute('href').substring(1);
            const targetElement = document.getElementById(targetId);
            
            if (targetElement) {
                e.preventDefault();
                const headerOffset = 80;
                const elementPosition = targetElement.getBoundingClientRect().top;
                const offsetPosition = elementPosition + window.pageYOffset - headerOffset;
                
                window.scrollTo({
                    top: offsetPosition,
                    behavior: 'smooth'
                });
            }
        });
    });
    
    // 添加页面加载动画
    document.body.classList.add('page-loaded');
}