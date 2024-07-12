document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll("video").forEach(function (video) {
    const videoSrc = video.getAttribute("data-source");
    if (hls.isSupported() && videoSrc) {
      var hls = new hls();
      hls.loadSource(videoSrc);
      hls.attachMedia(video);
    }
  });
});
