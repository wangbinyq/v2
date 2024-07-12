document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll("video").forEach(function (video) {
    const videoSrc = video.getAttribute("data-source");
    if (Hls.isSupported() && videoSrc) {
      var hls = new Hls();
      hls.loadSource(videoSrc);
      hls.attachMedia(video);
    }
  });
});
