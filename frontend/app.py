import streamlit as st
import requests
import os
import re
from datetime import datetime

PAGE_TITLE = os.environ.get('PAGE_TITLE', 'Web Chat Bot demo')
PAGE_HEADER = os.environ.get('PAGE_HEADER', 'Made with Streamlit and Parakeet')
PAGE_ICON = os.environ.get('PAGE_ICON', '🚀')

# Configuration of the Streamlit page
st.set_page_config(page_title=PAGE_TITLE, page_icon=PAGE_ICON)

# Hide the Deploy button
st.markdown("""
    <style>
    .stDeployButton {
        visibility: hidden;
    }
    .step-separator {
        border: none;
        border-top: 2px solid #4CAF50;
        margin: 10px 0;
    }
    .thinking-box {
        background-color: #f0f2f6;
        padding: 10px;
        border-radius: 5px;
        border-left: 4px solid #4CAF50;
        margin: 10px 0;
    }
    .step-indicator {
        background-color: #e8f4fd;
        padding: 8px;
        border-radius: 3px;
        border-left: 3px solid #2196F3;
        margin: 5px 0;
    }
    </style>
    """, unsafe_allow_html=True)

# Initialisation of the session state
if "messages" not in st.session_state:
    st.session_state.messages = []

# Handle the reset of the input key
if "input_key" not in st.session_state:
    st.session_state.input_key = 0

# Backend URL
BACKEND_SERVICE_URL = os.environ.get('BACKEND_SERVICE_URL', 'http://backend:5050')

def parse_and_render_content(content):
    """Parse content and render HTML elements appropriately"""
    
    # Liste des patterns HTML à supporter
    html_patterns = [
        # Séparateurs horizontaux
        (r'<hr>', '<hr class="step-separator">'),
        (r'<hr/>', '<hr class="step-separator">'),
        (r'<hr />', '<hr class="step-separator">'),
        
        # Messages de réflexion avec style spécial
        (r'<hr>🤖 ([^<]+)<hr>', r'<div class="thinking-box">🤖 \1</div>'),
        (r'🤖 ([^<\n]+)', r'<div class="thinking-box">🤖 \1</div>'),
        
        # Étapes avec numérotation
        (r'<step>([^<]+)</step>', r'<div class="step-indicator">📋 \1</div>'),
        (r'<info>([^<]+)</info>', r'<div class="step-indicator">ℹ️ \1</div>'),
        (r'<warning>([^<]+)</warning>', r'<div style="background-color: #fff3cd; padding: 8px; border-radius: 3px; border-left: 3px solid #ffc107; margin: 5px 0;">⚠️ \1</div>'),
        (r'<error>([^<]+)</error>', r'<div style="background-color: #f8d7da; padding: 8px; border-radius: 3px; border-left: 3px solid #dc3545; margin: 5px 0;">❌ \1</div>'),
        (r'<success>([^<]+)</success>', r'<div style="background-color: #d4edda; padding: 8px; border-radius: 3px; border-left: 3px solid #28a745; margin: 5px 0;">✅ \1</div>'),
        
        # Formatting basique
        (r'\*\*([^*]+)\*\*', r'<strong>\1</strong>'),
        (r'\*([^*]+)\*', r'<em>\1</em>'),
    ]
    
    # Séparer le contenu en parties HTML et non-HTML
    parts = []
    current_pos = 0
    
    # Chercher les balises HTML simples
    html_regex = r'<(hr|step|info|warning|error|success)[^>]*>.*?</?\1>?|<hr\s*/?>'
    
    for match in re.finditer(html_regex, content, re.IGNORECASE | re.DOTALL):
        # Ajouter le texte avant la balise HTML
        if match.start() > current_pos:
            text_part = content[current_pos:match.start()]
            if text_part.strip():
                parts.append(('text', text_part))
        
        # Ajouter la partie HTML
        html_part = match.group()
        parts.append(('html', html_part))
        current_pos = match.end()
    
    # Ajouter le reste du texte
    if current_pos < len(content):
        remaining_text = content[current_pos:]
        if remaining_text.strip():
            parts.append(('text', remaining_text))
    
    # Si aucune balise HTML trouvée, traiter tout comme du texte
    if not parts:
        parts = [('text', content)]
    
    return parts

def render_mixed_content(content):
    """Render content with mixed text and HTML"""
    parts = parse_and_render_content(content)
    
    for part_type, part_content in parts:
        if part_type == 'html':
            # Appliquer les transformations HTML
            processed_html = part_content
            html_patterns = [
                (r'<hr>', '<hr class="step-separator">'),
                (r'<hr/>', '<hr class="step-separator">'),
                (r'<hr />', '<hr class="step-separator">'),
                (r'<hr>🤖 ([^<]+)<hr>', r'<div class="thinking-box">🤖 \1</div>'),
                (r'<step>([^<]+)</step>', r'<div class="step-indicator">📋 \1</div>'),
                (r'<info>([^<]+)</info>', r'<div class="step-indicator">ℹ️ \1</div>'),
                (r'<warning>([^<]+)</warning>', r'<div style="background-color: #fff3cd; padding: 8px; border-radius: 3px; border-left: 3px solid #ffc107; margin: 5px 0;">⚠️ \1</div>'),
                (r'<error>([^<]+)</error>', r'<div style="background-color: #f8d7da; padding: 8px; border-radius: 3px; border-left: 3px solid #dc3545; margin: 5px 0;">❌ \1</div>'),
                (r'<success>([^<]+)</success>', r'<div style="background-color: #d4edda; padding: 8px; border-radius: 3px; border-left: 3px solid #28a745; margin: 5px 0;">✅ \1</div>'),
            ]
            
            for pattern, replacement in html_patterns:
                processed_html = re.sub(pattern, replacement, processed_html, flags=re.IGNORECASE)
            
            # Afficher le HTML
            st.markdown(processed_html, unsafe_allow_html=True)
        else:
            # Afficher le texte normal avec markdown
            if part_content.strip():
                st.markdown(part_content)

def stream_response(message):
    """Stream the message response from the backend"""
    try:
        with requests.post(
            BACKEND_SERVICE_URL+"/chat",
            json={"message": message},
            headers={"Content-Type": "application/json"},
            stream=True
        ) as response:
            # Create a placeholder for the streaming response
            response_placeholder = st.empty()
            full_response = ""
            
            # Stream the response chunks
            for chunk in response.iter_content(chunk_size=1024, decode_unicode=True):
                if chunk:
                    chunk_text = chunk.decode('utf-8') if isinstance(chunk, bytes) else chunk
                    full_response += chunk_text
                    
                    # Update the placeholder with the accumulated response using mixed content rendering
                    with response_placeholder.container():
                        render_mixed_content(full_response)
            
            return full_response
    except requests.exceptions.RequestException as e:
        error_msg = f"😡 Connection error: {str(e)}"
        st.error(error_msg)
        return error_msg

def increment_input_key():
    """Increment the input key to reset the input field"""
    st.session_state.input_key += 1

# Page title
st.header(PAGE_HEADER)

# Form to send a message
with st.form(key=f"message_form_{st.session_state.input_key}"):
    message = st.text_area("📝 Your message:", key=f"input_{st.session_state.input_key}", height=150)
    col1, col2, col3, col4, col5, col6 = st.columns(6)
    with col1:
        submit_button = st.form_submit_button(label="Send...")
    with col6:
        cancel_button = st.form_submit_button(label="Cancel", type="secondary")

# Handle the message submission
if submit_button and message and len(message.strip()) > 0:
    # Add the message to the history
    st.session_state.messages.append({
        "role": "user",
        "content": message,
        "time": datetime.now()
    })
    
    # Stream the response from the backend
    response = stream_response(message)
    
    # Add the response to the history
    st.session_state.messages.append({
        "role": "assistant",
        "content": response,
        "time": datetime.now()
    })
    
    # Reset the input field
    increment_input_key()
    st.rerun()

# Handle the message submission and cancellation
if cancel_button:
    try:
        response = requests.delete(f"{BACKEND_SERVICE_URL}/cancel")
        if response.status_code == 200:
            st.success("Request cancelled successfully")
        else:
            st.error("Failed to cancel request")
    except requests.exceptions.RequestException as e:
        st.error(f"Error cancelling request: {str(e)}")

# Display the messages history
st.write("### Messages history")
for msg in reversed(st.session_state.messages):
    with st.container():
        if msg["role"] == "user":
            st.info(f"🤓 You ({msg['time'].strftime('%H:%M')})")
            st.write(msg["content"])
        else:
            st.success(f"🤖 Assistant ({msg['time'].strftime('%H:%M')})")
            # Utiliser le rendu mixte pour l'historique aussi
            render_mixed_content(msg["content"])
            