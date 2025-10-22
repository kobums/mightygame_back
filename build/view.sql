drop view receptionlist_vw;
create view receptionlist_vw as select reception_tb.*, u_name as re_name, u_ssn as re_ssn from user_tb, reception_tb where u_id = re_user;
