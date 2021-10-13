<template>
  <div class="custom-foot">
    <div class="foot-main">
      <div class="copyright">Copyright &copy; 2021 BlockPI labs. All rights reserved.</div>
      <div class="foot-right">
        <div class="divide"></div>
        <a class="contact" href="https://twitter.com/RealBlockPI" target="_blank">
          <i class="iconfont icon-twitter"></i>
        </a>
        <a class="contact" href="https://t.me/blockpidaily" target="_blank">
          <i class="iconfont icon-telegram"></i>
        </a>
        <div class="divide"></div>
        <a class="mail" href="mailto:market@blockpi.io" target="_blank">
          <i class="iconfont icon-email"></i> <span>market@blockpi.io</span>
        </a>
      </div>
    </div>
  </div>
</template>

<script>
module.exports = {
  name: 'CustomFoot'
}
</script>

<style>
.custom-foot {
  height: 208px;
  width: 100%;
}
.foot-main {
  max-width: 1220px;
  padding: 0 10px;
  margin: 0px auto;
  height: 100%;
  background: url('../img/footer.png') no-repeat center bottom;
  background-size: cover;
  display: flex;
  align-items: center;
  justify-content: center;
}
.foot-main .copyright {
  position: relative;
  font-size: 12px;
  color: #5d687c;
}
.foot-main .copyright::before {
  position: absolute;
  bottom: -8px;
  left: 0;
  width: 80px;
  height: 1px;
  content: '';
  background: linear-gradient(to right,#9945ff, #00d18c) no-repeat left bottom;
}
.foot-main .foot-right {
  display: flex;
  align-items: center;
}
.foot-main .divide {
  width: 1px;
  height: 26px;
	background-color: #5d687c;
  margin: 0 30px;
}
.foot-main .contact .iconfont{
  cursor: pointer;
  font-size: 20px;
  margin: 0 10px;
  color: #00d18c;
}
.foot-main .mail{
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 16px;
  color: #5d687c;
  vertical-align: middle;
}
.foot-main .mail .iconfont {
  font-size: 20px;
  margin-right: 10px;
}
@media screen and (max-width: 768px) {
  .foot-main {
    flex-direction: column;
  }
  .foot-main .copyright {
    margin-bottom: 60px;
  }
}
</style>