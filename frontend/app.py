import streamlit as st
import requests
import os
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
    </style>
    """, unsafe_allow_html=True)

# Initialisation of the session state
if "messages" not in st.session_state:
    st.session_state.messages = []

# Handle the reset of the input key
if "input_key" not in st.session_state:
    st.session_state.input_key = 0

# Backend URL (the nodejs server)
#BACKEND_SERVICE_URL = "http://backend:5050"

BACKEND_SERVICE_URL = os.environ.get('BACKEND_SERVICE_URL', 'http://backend:5050')

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
                    # Update the placeholder with the accumulated response
                    response_placeholder.markdown(full_response)
            
            return full_response
    except requests.exceptions.RequestException as e:
        error_msg = f"😡 Connection error: {str(e)}"
        st.error(error_msg)
        return error_msg

def increment_input_key():
    """Increment the input key to reset the input field"""
    st.session_state.input_key += 1

# Page title
#st.title(PAGE_TITLE)
st.header(PAGE_HEADER)

# Form to send a message
with st.form(key=f"message_form_{st.session_state.input_key}"):
    #message = st.text_input("📝 Your message:", key=f"input_{st.session_state.input_key}")
    message = st.text_area("📝 Your message:", key=f"input_{st.session_state.input_key}", height=150)
    #submit_button = st.form_submit_button(label="Send...")
    #cancel_button = st.form_submit_button(label="Cancel", type="secondary")
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
            st.write(msg["content"])

