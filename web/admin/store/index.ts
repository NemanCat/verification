import { GetterTree, MutationTree } from 'vuex'

export const state = () => ({
  title: '',
  sidebar_collapsed:false,
  csfr_security_token: '',
  userid: '0',
  username: '',
  fetchError: {
    fetchErrorFlag: false,
    fetchErrorCode: 0,
    fetchErrorMessage: '',
  }
})


export type RootState = ReturnType<typeof state>

export const getters: GetterTree<RootState, RootState> = {
  title: state => state.title,
  sidebar_collapsed: state => state.sidebar_collapsed,
  csfr_security_token: state => state.csfr_security_token,
  userid: state => state.userid,
  username: state => state.username,
  fetchError: state => state.fetchError,
}

export const mutations: MutationTree<RootState> = {
  initializeStore: (state) => {
    if (process.browser) {
      const saved_sidebar_collapsed = window.sessionStorage.getItem('sidebar_collapsed');
      if(saved_sidebar_collapsed) {        
          state.sidebar_collapsed = (saved_sidebar_collapsed == 'true');       
      }  else {
        window.sessionStorage.setItem('sidebar_collapsed','false')
      }     
      const csfr_security_token = window.sessionStorage.getItem('csfr_security_token');
      if(csfr_security_token) {        
          state.csfr_security_token = csfr_security_token;       
      }  else {
        window.sessionStorage.setItem('csfr_security_token','')
      } 
      const userid = window.sessionStorage.getItem('userid');
      
      if(userid) {        
          state.userid = userid;       
      }  else {
        window.sessionStorage.setItem('userid','0')
      }     
      const username = window.sessionStorage.getItem('username');
      if(username) {        
          state.username = username;       
      }  else {
        window.sessionStorage.setItem('username','')
      }
    } 
  },
  CHANGE_TITLE: (state, newTitle: string) => (state.title = newTitle),
  CHANGE_SIDEBAR_COLLAPSED: (state,newSidebarCollapsed: boolean) => {
    state.sidebar_collapsed = newSidebarCollapsed;
    if (process.browser) {
      window.sessionStorage.setItem('sidebar_collapsed',String(state.sidebar_collapsed));
    }      
  },  
  CHANGE_CSFR_SECURITY_TOKEN: (state, newCsfrSecurityToken: string) => {
    state.csfr_security_token = newCsfrSecurityToken;
    if (process.browser) {
      window.sessionStorage.setItem('csfr_security_token',state.csfr_security_token);
    }      
  },
  CHANGE_USERID: (state, newUserid: string) => {
    state.userid = newUserid;
    if (process.browser) {
      window.sessionStorage.setItem('userid',state.userid);
    }      
  },
  CHANGE_USERNAME: (state, newUsername: string) => {
    state.username = newUsername;
    if (process.browser) {
      window.sessionStorage.setItem('username',state.username);
    }      
  }, 
  CHANGE_FETCHERROR: (state, newFetchError) => {
    state.fetchError = newFetchError;
  },
}


  