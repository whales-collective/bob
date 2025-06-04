import streamlit as st
import requests
import os
import re
from datetime import datetime

#PAGE_TITLE = os.environ.get('PAGE_TITLE', 'Web Chat Bot demo')
PAGE_HEADER = os.environ.get('PAGE_HEADER', 'Made with Streamlit and LangChainJS')

PAGE_ICON = os.environ.get('PAGE_ICON', '🚀')

LLM_CHAT= os.environ.get('MODEL_RUNNER_CHAT_MODEL_BOB', '')
LLM_EMBEDDINGS= os.environ.get('MODEL_RUNNER_EMBEDDING_MODEL', '')

# Configuration of the Streamlit page
#st.set_page_config(page_title=PAGE_TITLE, page_icon=PAGE_ICON)

# Hide the Deploy button and adjust header position
st.markdown("""
    <style>
    .stDeployButton {
        visibility: hidden;
    }
    /* Ensure fonts render correctly for monospace content */
    code, pre {
        font-family: 'Courier New', Courier, monospace;
    }
    /* Reduce top margin/padding for the main content */
    .main .block-container {
        padding-top: 1rem;
        padding-bottom: 10rem;
        margin-top: 0.7rem;
    }
    /* Adjust header spacing */
    h2 {
        text-align: center !important;
        margin-top: 0rem !important;
        padding-top: 0rem !important;
    }
    h1, h3 {
        margin-top: 0rem !important;
        padding-top: 0rem !important;
    }
    /* Force buttons to align right */
    .stForm > div:last-child {
        display: flex !important;
        justify-content: flex-end !important;
        gap: 10px !important;
    }
    .stForm > div:last-child > div {
        flex: none !important;
    }
    
    /* GitHub-style labels */
    .github-label {
        display: inline-block;
        padding: 2px 6px;
        margin: 2px;
        border-radius: 12px;
        font-size: 12px;
        font-weight: 500;
        color: white;
        text-decoration: none;
        white-space: nowrap;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    }
    
    /* Predefined label colors (GitHub-inspired) */
    .label-red { background-color: #d73a49; }
    .label-green { background-color: #28a745; }
    .label-blue { background-color: #0366d6; }
    .label-yellow { background-color: #ffd33d; color: #24292e; }
    .label-orange { background-color: #f66a0a; }
    .label-purple { background-color: #6f42c1; }
    .label-pink { background-color: #ea4aaa; }
    .label-gray { background-color: #6a737d; }
    .label-grey { background-color: #6a737d; }
    .label-black { background-color: #24292e; }
    .label-white { background-color: #f6f8fa; color: #24292e; border: 1px solid #d1d5da; }
    
    /* Status-specific colors */
    .label-success { background-color: #28a745; }
    .label-error { background-color: #d73a49; }
    .label-warning { background-color: #ffd33d; color: #24292e; }
    .label-info { background-color: #0366d6; }
    .label-step { background-color: #6f42c1; }
    
    /* Additional semantic colors */
    .label-bug { background-color: #d73a49; }
    .label-feature { background-color: #a2eeef; color: #24292e; }
    .label-documentation { background-color: #0075ca; }
    .label-enhancement { background-color: #a2eeef; color: #24292e; }
    .label-question { background-color: #d876e3; }
    .label-wontfix { background-color: #ffffff; color: #24292e; border: 1px solid #d1d5da; }
    .label-duplicate { background-color: #cfd3d7; color: #24292e; }
    .label-invalid { background-color: #e4e669; color: #24292e; }
    
    /* Default fallback */
    .label-default { background-color: #586069; }
    </style>
    """, unsafe_allow_html=True)

# Initialisation of the session state
if "messages" not in st.session_state:
    st.session_state.messages = []

# Handle the reset of the input key
if "input_key" not in st.session_state:
    st.session_state.input_key = 0

# Initialize session ID in session state
if "session_id" not in st.session_state:
    st.session_state.session_id = "default"


# Backend URL (the nodejs server)
BACKEND_SERVICE_URL = os.environ.get('BACKEND_SERVICE_URL', 'http://backend:5050')

def convert_tags_to_labels(text):
    """Convert HTML-like tags to GitHub-style labels"""
    # Pattern to match tags like <red>content</red>, <step>content</step>, etc.
    pattern = r'<(\w+)>(.*?)</\1>'
    
    def replace_tag(match):
        tag_name = match.group(1).lower()
        content = match.group(2)
        
        # Create a GitHub-style label
        label_class = f"label-{tag_name}"
        return f'<span class="github-label {label_class}">{content}</span>'
    
    # Replace all matching tags
    result = re.sub(pattern, replace_tag, text)
    return result

def stream_response(message, session_id):
    """Stream the message response from the backend"""
    try:
        with requests.post(
            BACKEND_SERVICE_URL+"/chat",
            json={"message": message, "sessionId": session_id},
            headers={"Content-Type": "application/json"},
            stream=True
        ) as response:
            # Create a placeholder for the streaming response
            response_placeholder = st.empty()
            full_response = ""
            
            # Stream the response chunks
            for chunk in response.iter_content(chunk_size=1024):
                if chunk:
                    try:
                        # Decode using utf-8 with error handling
                        chunk_text = chunk.decode('utf-8', errors='replace')
                        full_response += chunk_text
                        
                        # Convert tags to labels and process content
                        processed_response = convert_tags_to_labels(full_response)
                        
                        # For content that might contain tree-view or other special characters
                        if "```" in processed_response:
                            # Process code blocks to ensure tree structures are formatted as raw
                            formatted_response = process_code_blocks(processed_response)
                            response_placeholder.markdown(formatted_response, unsafe_allow_html=True)
                        else:
                            response_placeholder.markdown(processed_response, unsafe_allow_html=True)
                    except UnicodeDecodeError:
                        # If there's still a decode error, replace problematic characters
                        st.warning("Encountered encoding issues with response")
                        chunk_text = chunk.decode('utf-8', errors='replace')
                        full_response += chunk_text
                        processed_response = convert_tags_to_labels(full_response)
                        response_placeholder.markdown(processed_response, unsafe_allow_html=True)
            
            return full_response
    except requests.exceptions.RequestException as e:
        error_msg = f"😡 Connection error: {str(e)}"
        st.error(error_msg)
        return error_msg

def process_code_blocks(text):
    """Process code blocks to properly display tree-view and other special characters"""
    # Split by code block markers
    parts = text.split("```")
    
    result = []
    for i, part in enumerate(parts):
        if i % 2 == 0:  # This is regular text
            result.append(part)
        else:  # This is code
            # Get the language identifier if present
            lines = part.strip().split('\n', 1)
            lang = lines[0].strip() if len(lines) > 1 else ""
            code_content = lines[1] if len(lines) > 1 else lines[0]
            
            # Check if it's a tree structure (simple heuristic)
            if any(char in code_content for char in ['│', '├', '└', '─', '┬', '┤']):
                # For tree structures, use raw formatting
                if lang != "raw":
                    result.append(f"```raw\n{code_content}\n```")
                else:
                    result.append(f"```{lang}\n{code_content}\n```")
            else:
                # Regular code, keep the original format
                result.append(f"```{part}```")
    
    return "".join(result)

def clear_conversation_history(session_id):
    """Clear the conversation history on the server"""
    try:
        response = requests.post(
            f"{BACKEND_SERVICE_URL}/clear-history",
            json={"sessionId": session_id},
            headers={"Content-Type": "application/json"}
        )
        if response.status_code == 200:
            st.session_state.messages = []  # Clear local messages too
            st.success("✨ Conversation history cleared!")
        else:
            st.error("Failed to clear conversation history")
    except requests.exceptions.RequestException as e:
        st.error(f"Error clearing history: {str(e)}")


def increment_input_key():
    """Increment the input key to reset the input field"""
    st.session_state.input_key += 1

# Page title - moved higher with reduced spacing
#st.title(PAGE_TITLE)
st.header(PAGE_HEADER)

# Session ID input
session_id = st.text_input(
    "🔑 Session ID:",
    value=st.session_state.session_id,
    help="Enter a unique session ID to maintain conversation context"
)
st.session_state.session_id = session_id

#models_info = st.text_input("Models", value=f"📕 {LLM_CHAT} - 🌐 {LLM_EMBEDDINGS}")

# Form to send a message
with st.form(key=f"message_form_{st.session_state.input_key}"):
    message = st.text_area("📝 Your message:", key=f"input_{st.session_state.input_key}", height=80)
    
    # Utilisation de colonnes pour pousser les boutons à droite
    col1, col2, col3, col4 = st.columns([10, 6, 3, 3])
    with col1:
        st.empty()  # Colonne vide pour pousser les boutons à droite
    with col2:
        st.empty()
    with col3:
        submit_button = st.form_submit_button(label="Send...")
    with col4:
        cancel_button = st.form_submit_button(label="Cancel", type="secondary")
        

# Handle the message submission
if submit_button and message and len(message.strip()) > 0:
    # Add the message to the history
    st.session_state.messages.append({
        "role": "user",
        "content": message,
        "time": datetime.now(),
        "session_id": st.session_state.session_id
    })
    
    # Stream the response from the backend
    response = stream_response(message, st.session_state.session_id)
    
    # Add the response to the history
    st.session_state.messages.append({
        "role": "assistant",
        "content": response,
        "time": datetime.now(),
        "session_id": st.session_state.session_id
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
#st.write("### Messages history")
for msg in reversed(st.session_state.messages):
    with st.container():
        if msg["role"] == "user":
            st.info(f"🤓 You ({msg['time'].strftime('%H:%M')}) - Session: {msg['session_id']}")
            st.write(msg["content"])
        else:
            st.success(f"🤖 Assistant ({msg['time'].strftime('%H:%M')}) - Session: {msg['session_id']}")
            # Process the message content to handle special formatting and labels
            processed_content = convert_tags_to_labels(msg["content"])
            
            if "```" in processed_content:
                formatted_content = process_code_blocks(processed_content)
                st.markdown(formatted_content, unsafe_allow_html=True)
            else:
                st.markdown(processed_content, unsafe_allow_html=True)