FROM python:3.9-slim

WORKDIR /app

# Install required packages
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY app.py .
COPY templates/ templates/

# Set environment variable for the flag
ENV FLAG="CTF{j1nj4_sst1_vuln3r4bl3_t3mpl4t3s}"

# Create a non-root user
RUN useradd -m ctfuser
USER ctfuser

# Expose the port
EXPOSE 5000

# Run the application
CMD ["python", "app.py"]
