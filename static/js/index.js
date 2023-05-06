const mediaSource = new MediaSource();

async function getVideo() {
  const res = await fetch("http://localhost:5001/api/videos", {
    headers: { "Content-Type": "application/json" },
    method: "POST",
    body: JSON.stringify({ name: "sheesh.mp4" }),
  });
  const data = await res.json();
  console.log(data);
  document.getElementById("videoPlayer").src = data.videoUrl;
}
