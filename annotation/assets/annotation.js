(function() {
    'use strict';
    
    const API_ENDPOINT = window.ANNOTATION_API || '/api/annotations';
    let annotationMode = false;
    let annotations = [];
    
    // Initialize annotation system
    function init() {
        createAnnotationUI();
        loadAnnotations();
        setupEventListeners();
        console.log('Annotation system initialized');
    }
    
    // Create UI overlay
    function createAnnotationUI() {
        const toolbar = document.createElement('div');
        toolbar.id = 'annotation-toolbar';
        toolbar.innerHTML = `
            <button id="toggle-annotation-mode">📝 Annotate</button>
            <button id="show-annotations">👁️ View (<span id="annotation-count">0</span>)</button>
        `;
        document.body.appendChild(toolbar);
    }
    
    // Setup event listeners
    function setupEventListeners() {
        document.getElementById('toggle-annotation-mode').addEventListener('click', toggleAnnotationMode);
        document.getElementById('show-annotations').addEventListener('click', showAnnotationList);
        
        document.addEventListener('click', handleClick, true);
    }
    
    // Toggle annotation mode
    function toggleAnnotationMode() {
        annotationMode = !annotationMode;
        document.body.classList.toggle('annotation-mode-active', annotationMode);
        
        const btn = document.getElementById('toggle-annotation-mode');
        btn.textContent = annotationMode ? '✓ Annotating' : '📝 Annotate';
        btn.style.background = annotationMode ? '#4CAF50' : '#2196F3';
    }
    
    // Handle element click in annotation mode
    function handleClick(e) {
        if (!annotationMode) return;
        
        const target = e.target;
        
        // Ignore clicks on annotation UI elements
        if (target.closest('#annotation-toolbar') || 
            target.closest('.annotation-marker') ||
            target.closest('#annotation-form-modal') ||
            target.closest('.annotation-form')) {
            return;
        }
        
        e.preventDefault();
        e.stopPropagation();
        
        showAnnotationForm(target, e.clientX, e.clientY);
    }
    
    // Show annotation form
    function showAnnotationForm(element, x, y) {
        const existingForm = document.getElementById('annotation-form-modal');
        if (existingForm) existingForm.remove();
        
        const modal = document.createElement('div');
        modal.id = 'annotation-form-modal';
        modal.innerHTML = `
            <div class="annotation-form">
                <h3>Add Annotation</h3>
                <input type="text" id="author-input" placeholder="Your name" />
                <textarea id="comment-input" placeholder="Your feedback..."></textarea>
                <div class="form-actions">
                    <button id="save-annotation">Save</button>
                    <button id="cancel-annotation">Cancel</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
        
        // Close on backdrop click
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
        
        // Close on Escape key
        const escapeHandler = (e) => {
            if (e.key === 'Escape') {
                modal.remove();
                document.removeEventListener('keydown', escapeHandler);
            }
        };
        document.addEventListener('keydown', escapeHandler);
        
        document.getElementById('save-annotation').onclick = () => saveAnnotation(element, x, y);
        document.getElementById('cancel-annotation').onclick = () => modal.remove();
        document.getElementById('comment-input').focus();
    }
    
    // Save annotation
    async function saveAnnotation(element, x, y) {
        const author = document.getElementById('author-input').value.trim();
        const comment = document.getElementById('comment-input').value.trim();
        
        if (!comment) {
            alert('Please enter a comment');
            return;
        }
        
        const annotation = {
            url: window.location.pathname,
            selector: getSelector(element),
            comment: comment,
            author: author || 'Anonymous',
            position: { x, y }
        };
        
        try {
            const response = await fetch(API_ENDPOINT, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(annotation)
            });
            
            if (response.ok) {
                const saved = await response.json();
                annotations.push(saved);
                renderAnnotationMarker(saved);
                updateAnnotationCount();
                document.getElementById('annotation-form-modal').remove();
                toggleAnnotationMode();
            }
        } catch (error) {
            console.error('Failed to save annotation:', error);
            alert('Failed to save annotation');
        }
    }
    
    // Load existing annotations
    async function loadAnnotations() {
        try {
            const response = await fetch(`${API_ENDPOINT}?url=${encodeURIComponent(window.location.pathname)}`);
            const data = await response.json();
            annotations = data.annotations || [];
            renderAnnotations();
            updateAnnotationCount();
        } catch (error) {
            console.error('Failed to load annotations:', error);
        }
    }
    
    // Render all annotations
    function renderAnnotations() {
        annotations.forEach(renderAnnotationMarker);
    }
    
    // Render single annotation marker
    function renderAnnotationMarker(annotation) {
        const marker = document.createElement('div');
        marker.className = 'annotation-marker';
        marker.style.left = annotation.position.x + 'px';
        marker.style.top = annotation.position.y + 'px';
        marker.innerHTML = '💬';
        marker.title = `${annotation.author}: ${annotation.comment}`;
        
        marker.onclick = () => showAnnotationDetail(annotation);
        document.body.appendChild(marker);
    }
    
    // Show annotation detail
    function showAnnotationDetail(annotation) {
        const modal = document.createElement('div');
        modal.id = 'annotation-form-modal';
        modal.innerHTML = `
            <div class="annotation-form">
                <h3>Annotation</h3>
                <p><strong>${annotation.author}</strong></p>
                <p>${annotation.comment}</p>
                <div class="form-actions">
                    <button id="close-annotation">Close</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
        
        // Close on backdrop click
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
        
        // Close on Escape key
        const escapeHandler = (e) => {
            if (e.key === 'Escape') {
                modal.remove();
                document.removeEventListener('keydown', escapeHandler);
            }
        };
        document.addEventListener('keydown', escapeHandler);
        
        document.getElementById('close-annotation').onclick = () => modal.remove();
    }
    
    // Show list of all annotations
    function showAnnotationList() {
        const modal = document.createElement('div');
        modal.id = 'annotation-form-modal';
        
        let html = '<div class="annotation-form"><h3>All Annotations</h3>';
        if (annotations.length === 0) {
            html += '<p>No annotations yet</p>';
        } else {
            html += '<ul class="annotation-list">';
            annotations.forEach(ann => {
                html += `<li><strong>${ann.author}:</strong> ${ann.comment}</li>`;
            });
            html += '</ul>';
        }
        html += '<button id="close-list">Close</button></div>';
        
        modal.innerHTML = html;
        document.body.appendChild(modal);
        
        // Close on backdrop click
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
        
        // Close on Escape key
        const escapeHandler = (e) => {
            if (e.key === 'Escape') {
                modal.remove();
                document.removeEventListener('keydown', escapeHandler);
            }
        };
        document.addEventListener('keydown', escapeHandler);
        
        document.getElementById('close-list').onclick = () => modal.remove();
    }
    
    // Update annotation count
    function updateAnnotationCount() {
        const countEl = document.getElementById('annotation-count');
        if (countEl) {
            countEl.textContent = annotations.length;
        }
    }
    
    // Get CSS selector for element
    function getSelector(element) {
        if (element.id) return `#${element.id}`;
        if (element.className) {
            const classes = element.className.split(' ').filter(c => c).join('.');
            return element.tagName.toLowerCase() + '.' + classes;
        }
        return element.tagName.toLowerCase();
    }
    
    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
