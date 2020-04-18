import { ReactComponent as EditIcon } from 'assets/edit-pencil.svg';
import { ReactComponent as NoUser } from 'assets/no-user-icon.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloModal } from 'Components/Common/SoloModal';
import { Portal } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/portal_pb';
import React from 'react';
import useSWR from 'swr';
import { userApi } from '../api';
import { AddUserModal } from './AddUserModal';
import { User } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/user_pb';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { Loading } from 'Components/Common/DisplayOnly/Loading';

type PortalUsersTabProps = {
  portal: Portal.AsObject;
};

export const PortalUsersTab = ({ portal }: PortalUsersTabProps) => {
  const { data: usersList, error: usersError } = useSWR(
    `listUsers${portal.metadata?.name}${portal.metadata?.namespace}`,
    () =>
      userApi.listUsers({
        portalsList: [
          { name: portal.metadata!.name, namespace: portal.metadata!.namespace }
        ],
        apiDocsList: [],
        groupsList: []
      })
  );

  const [userSearchTerm, setUserSearchTerm] = React.useState('');
  const [showCreateUserModal, setShowCreateUserModal] = React.useState(false);

  const [showConfirmUserDelete, setShowConfirmUserDelete] = React.useState(
    false
  );
  const [userToDelete, setUserToDelete] = React.useState<User.AsObject>();
  const [filteredUsers, setFilteredUsers] = React.useState<User.AsObject[]>();

  React.useEffect(() => {
    if (userSearchTerm !== '' && !!usersList) {
      setFilteredUsers(
        usersList.filter(user =>
          user.spec?.username?.toLowerCase().includes(userSearchTerm)
        )
      );
    } else {
      setFilteredUsers(undefined);
    }
  }, [userSearchTerm]);

  const attemptDeleteUser = (user: User.AsObject) => {
    setShowConfirmUserDelete(true);
    setUserToDelete(user);
  };

  const deleteUser = async () => {
    await userApi.deleteUser({
      name: userToDelete?.metadata?.name!,
      namespace: userToDelete?.metadata?.namespace!
    });
    closeConfirmModal();
  };

  const closeConfirmModal = () => {
    setShowConfirmUserDelete(false);
  };
  if (!usersList) {
    return <Loading center>Loading...</Loading>;
  }
  return (
    <div className='relative flex flex-col p-4 border border-gray-300 rounded-lg'>
      <span
        onClick={() => setShowCreateUserModal(true)}
        className='absolute top-0 right-0 flex items-center mt-2 mr-2 text-green-400 cursor-pointer hover:text-green-300'>
        <GreenPlus className='mr-1 fill-current' />
        <span className='text-gray-700'> Add/Remove a User</span>
      </span>
      <div className='w-1/3 m-4'>
        <SoloInput
          placeholder='Search by user name...'
          value={userSearchTerm}
          onChange={e => setUserSearchTerm(e.target.value)}
        />
      </div>
      <div className='flex flex-col'>
        <div className='py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8'>
          <div className='inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg'>
            <table className='min-w-full'>
              <thead className='bg-gray-300 '>
                <tr>
                  <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                    User Name
                  </th>
                  <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                    Email
                  </th>

                  <th className='px-6 py-3 text-sm font-medium leading-4 tracking-wider text-left text-gray-800 capitalize border-b border-gray-200 bg-gray-50'>
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className='bg-white'>
                {(!!filteredUsers ? filteredUsers : usersList)
                  .sort((a, b) =>
                    a.metadata?.name === b.metadata?.name
                      ? 0
                      : a.metadata!.name > b.metadata!.name
                      ? 1
                      : -1
                  )
                  .map(user => {
                    return (
                      <tr key={user.metadata?.uid}>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {user?.spec?.username}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 whitespace-no-wrap border-b border-gray-200'>
                          <div className='text-sm leading-5 text-gray-900'>
                            <span className='flex items-center capitalize'>
                              {user?.spec?.email}
                            </span>
                          </div>
                        </td>
                        <td className='px-6 py-4 text-sm font-medium leading-5 text-right whitespace-no-wrap border-b border-gray-200'>
                          <span className='flex items-center'>
                            <div className='flex items-center justify-center w-4 h-4 mr-3 text-gray-700 bg-gray-400 rounded-full cursor-pointer'>
                              <EditIcon className='w-2 h-3 fill-current' />
                            </div>
                            <div
                              className='flex items-center justify-center w-4 h-4 text-gray-700 bg-gray-400 rounded-full cursor-pointer'
                              onClick={() => attemptDeleteUser(user)}>
                              x
                            </div>
                          </span>
                        </td>
                      </tr>
                    );
                  })}
              </tbody>
            </table>
            {(!!filteredUsers ? filteredUsers : usersList).length === 0 && (
              <div className='w-full m-auto'>
                <div className='flex flex-col items-center justify-center w-full h-full py-4 mr-32 bg-white rounded-lg shadow-lg md:flex-row'>
                  <div className='mr-6'>
                    <NoUser />
                  </div>
                  <div className='flex flex-col h-full'>
                    <p className='h-auto my-6 text-lg font-medium text-gray-800 '>
                      There are no Users to display!{' '}
                    </p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
      <SoloModal visible={showCreateUserModal} width={750} noPadding={true}>
        <AddUserModal onClose={() => setShowCreateUserModal(false)} />
      </SoloModal>
      <ConfirmationModal
        visible={showConfirmUserDelete}
        confirmationTopic='delete this user'
        confirmText='Delete'
        goForIt={deleteUser}
        cancel={closeConfirmModal}
        isNegative={true}
      />
    </div>
  );
};
