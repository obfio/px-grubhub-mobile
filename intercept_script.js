// ! unpin
var TrustManagerImpl = Java.use('com.android.org.conscrypt.TrustManagerImpl');
TrustManagerImpl.verifyChain.implementation = function (untrustedChain, trustAnchorChain, host, clientAuth, ocspData, tlsSctData) {
    return untrustedChain;
};
Java.perform(function() {
    let okhttp3Pin = Java.use("okhttp3.CertificatePinner$Builder")
    okhttp3Pin["add"].implementation = function(pattern, pins) {
        return this
    }
})
// ! PX247
Java.perform(function(){
    let a = Java.use("ov0.a");
a["a"].implementation = function (i12, i13, i14, i15) {
    console.log('a is called' + ', ' + 'i12: ' + i12 + ', ' + 'i13: ' + i13 + ', ' + 'i14: ' + i14 + ', ' + 'i15: ' + i15);
    let ret = this.a(i12, i13, i14, i15);
    console.log('a ret value is ' + ret);
    return ret;
};
})
Java.perform(function(){
    let Boxing = Java.use("kotlin.coroutines.jvm.internal.Boxing");
Boxing["boxInt"].implementation = function (i12) {
    console.log('boxInt is called' + ', ' + 'i12: ' + i12);
    let ret = this.boxInt(i12);
    console.log('boxInt ret value is ' + ret);
    return ret;
};
})
