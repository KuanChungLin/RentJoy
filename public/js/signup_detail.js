const form = document.querySelector("#signinForm"); const formCheckInput = document.querySelector(".form-check-input");
const signupPasswordConfirm = document.querySelector(".signup-password-confirm");
const signupPassword = document.querySelector(".signup-password");
const phoneNumber = document.querySelector(".phonenumber");


form.addEventListener('submit', function (e) {
    if (!phoneNumber.value.match(/^09\d{8}$/)) {
        e.preventDefault();
        swal({
            title: "輸入的手機號碼須為09XXXXXXXX！",
            icon: "warning",
            button: true,
            dangerMode: true
        });
        return;
    }

    if (signupPasswordConfirm.value !== signupPassword.value) {
        e.preventDefault();
        swal({
            title: "再次輸入密碼的值與密碼不相符！",
            icon: "warning",
            button: true,
            dangerMode: true
        });
        signupPasswordConfirm.value = "";
        return;
    }

    if (!formCheckInput.checked) {
        e.preventDefault();
        swal({
            title: "請閱讀並同意服務條款！",
            icon: "warning",
            button: true,
            dangerMode: true
        });
        return;
    }
})