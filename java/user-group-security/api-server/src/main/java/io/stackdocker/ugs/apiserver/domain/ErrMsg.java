package io.stackdocker.ugs.apiserver.domain;

public class ErrMsg {

    private String error;

    public ErrMsg(String error) {
        this.error = error;
    }

    public String getError() {
        return error;
    }

    public void setError(String error) {
        this.error = error;
    }
}
