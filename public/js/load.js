var data;
var xhttp = new XMLHttpRequest();
xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
        data = JSON.parse(this.responseText);
    };
};
xhttp.open("GET", "/data/data.json", false);
xhttp.send();
var map = new AMap.Map('container', {
    resizeEnable: true,
    zoom: 4,
    center: [105, 33],
});
for (school in data) {
    studensStr = school + ":</br>";
    data[school].students.forEach(function(stu) {
        studensStr = studensStr + stu + "</br>";
    });
    marker = new AMap.Marker({
        position: data[school].pos,
        title: school,
        label: {
            offset: new AMap.Pixel(-24, 20),
            content: '<div class="info">' + studensStr + '</div>',
        }
    });
    map.add(marker);
}