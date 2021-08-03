var data;
var classIndex;
var datareq = new XMLHttpRequest();
var indexreq = new XMLHttpRequest();
indexreq.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
        classIndex = JSON.parse(this.responseText);
    };
};
indexreq.open("GET", "./data/index.json", false)
indexreq.send()
var map = new AMap.Map('container', {
    resizeEnable: true,
    zoom: 4,
    center: [105, 33],
});
ElementPlus.ElNotification({
    title: "书生21届毕业去向",
    message: '<p>技术支持：柏喵Sakura <a href="https://github.com/MscBaiMeow/sunsonMap" target="_blank"><i class="fa fa-github"></i></a></p><p>数据统计：特别贫穷的小C，鳝</p>',
    dangerouslyUseHTMLString: true,
});
var Main = {
    data() {
        return {
            chosenClass: 9,
            classIndex: classIndex,
        }
    },
    methods: {
        chooseClass(i) {
            loadClass(i);
            this.chosenClass = i;
            this.$notify({
                title: '成功',
                message: '成功切换到' + i + '班',
                type: 'success'
            });
        }

    }
}
const app = Vue.createApp(Main);
app.use(ElementPlus);
app.mount('#app');
loadClass(9)

function loadClass(classID) {
    map.clearMap();
    datareq.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            data = JSON.parse(this.responseText);
        };
    };
    datareq.open("GET", "./data/data" + classID + ".json", false);
    datareq.send();
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
}
